package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	goabi "github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	cltypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_consumer_benchmark"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_chain_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_consumer_performance_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_triggered_streams_lookup_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/perform_data_checker_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/streams_lookup_upkeep_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_transcoder"
)

var compatibleUtils = cltypes.MustGetABI(ac.AutomationCompatibleUtilsABI)
var registrarABI = cltypes.MustGetABI(registrar21.AutomationRegistrarABI)

type KeeperRegistrar interface {
	Address() string

	EncodeRegisterRequest(name string, email []byte, upkeepAddr string, gasLimit uint32, adminAddr string, checkData []byte, amount *big.Int, source uint8, senderAddr string, isLogTrigger bool, isMercury bool) ([]byte, error)

	Fund(ethAmount *big.Float) error
}

type UpkeepTranscoder interface {
	Address() string
}

type KeeperRegistry interface {
	Address() string
	Fund(ethAmount *big.Float) error
	SetConfig(config KeeperRegistrySettings, ocrConfig OCRv2Config) error
	SetConfigTypeSafe(ocrConfig OCRv2Config) error
	SetRegistrar(registrarAddr string) error
	AddUpkeepFunds(id *big.Int, amount *big.Int) error
	AddUpkeepFundsFromKey(id *big.Int, amount *big.Int, keyNum int) error
	GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error)
	GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error)
	SetKeepers(keepers []string, payees []string, ocrConfig OCRv2Config) error
	GetKeeperList(ctx context.Context) ([]string, error)
	RegisterUpkeep(target string, gasLimit uint32, admin string, checkData []byte) error
	CancelUpkeep(id *big.Int) error
	SetUpkeepGasLimit(id *big.Int, gas uint32) error
	ParseUpkeepPerformedLog(log *types.Log) (*UpkeepPerformedLog, error)
	ParseStaleUpkeepReportLog(log *types.Log) (*StaleUpkeepReportLog, error)
	ParseUpkeepIdFromRegisteredLog(log *types.Log) (*big.Int, error)
	Pause() error
	Migrate(upkeepIDs []*big.Int, destinationAddress common.Address) error
	SetMigrationPermissions(peerAddress common.Address, permission uint8) error
	PauseUpkeep(id *big.Int) error
	UnpauseUpkeep(id *big.Int) error
	UpdateCheckData(id *big.Int, newCheckData []byte) error
	SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) error
	SetUpkeepPrivilegeConfig(id *big.Int, privilegeConfig []byte) error
	SetUpkeepOffchainConfig(id *big.Int, offchainConfig []byte) error
	RegistryOwnerAddress() common.Address
	ChainModuleAddress() common.Address
	ReorgProtectionEnabled() bool
}

type KeeperConsumer interface {
	Address() string
	Counter(ctx context.Context) (*big.Int, error)
	Start() error
}

type UpkeepCounter interface {
	Address() string
	Fund(ethAmount *big.Float) error
	Counter(ctx context.Context) (*big.Int, error)
	SetSpread(testRange *big.Int, interval *big.Int) error
	Start() error
}

type UpkeepPerformCounterRestrictive interface {
	Address() string
	Fund(ethAmount *big.Float) error
	Counter(ctx context.Context) (*big.Int, error)
	SetSpread(testRange *big.Int, interval *big.Int) error
}

// KeeperConsumerPerformance is a keeper consumer contract that is more complicated than the typical consumer,
// it's intended to only be used for performance tests.
type KeeperConsumerPerformance interface {
	Address() string
	Fund(ethAmount *big.Float) error
	CheckEligible(ctx context.Context) (bool, error)
	GetUpkeepCount(ctx context.Context) (*big.Int, error)
	SetCheckGasToBurn(ctx context.Context, gas *big.Int) error
	SetPerformGasToBurn(ctx context.Context, gas *big.Int) error
}

// AutomationConsumerBenchmark is a keeper consumer contract that is more complicated than the typical consumer,
// it's intended to only be used for benchmark tests.
type AutomationConsumerBenchmark interface {
	Address() string
	Fund(ethAmount *big.Float) error
	CheckEligible(ctx context.Context, id *big.Int, _range *big.Int, firstEligibleBuffer *big.Int) (bool, error)
	GetUpkeepCount(ctx context.Context, id *big.Int) (*big.Int, error)
}

type KeeperPerformDataChecker interface {
	Address() string
	Counter(ctx context.Context) (*big.Int, error)
	SetExpectedData(ctx context.Context, expectedData []byte) error
}

type UpkeepPerformedLog struct {
	Id      *big.Int
	Success bool
	From    common.Address
}

type StaleUpkeepReportLog struct {
	Id *big.Int
}

// KeeperRegistryOpts opts to deploy keeper registry version
type KeeperRegistryOpts struct {
	RegistryVersion ethereum.KeeperRegistryVersion
	LinkAddr        string
	ETHFeedAddr     string
	GasFeedAddr     string
	TranscoderAddr  string
	RegistrarAddr   string
	Settings        KeeperRegistrySettings
}

// KeeperRegistrySettings represents the settings to fine tune keeper registry
type KeeperRegistrySettings struct {
	PaymentPremiumPPB    uint32   // payment premium rate oracles receive on top of being reimbursed for gas, measured in parts per billion
	FlatFeeMicroLINK     uint32   // flat fee charged for each upkeep
	BlockCountPerTurn    *big.Int // number of blocks each oracle has during their turn to perform upkeep before it will be the next keeper's turn to submit
	CheckGasLimit        uint32   // gas limit when checking for upkeep
	StalenessSeconds     *big.Int // number of seconds that is allowed for feed data to be stale before switching to the fallback pricing
	GasCeilingMultiplier uint16   // multiplier to apply to the fast gas feed price when calculating the payment ceiling for keepers
	MinUpkeepSpend       *big.Int // minimum spend required by an upkeep before they can withdraw funds
	MaxPerformGas        uint32   // max gas allowed for an upkeep within perform
	FallbackGasPrice     *big.Int // gas price used if the gas price feed is stale
	FallbackLinkPrice    *big.Int // LINK price used if the LINK price feed is stale
	MaxCheckDataSize     uint32
	MaxPerformDataSize   uint32
	MaxRevertDataSize    uint32
	RegistryVersion      ethereum.KeeperRegistryVersion
}

// KeeperRegistrarSettings represents settings for registrar contract
type KeeperRegistrarSettings struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint16
	RegistryAddr          string
	MinLinkJuels          *big.Int
}

// KeeperInfo keeper status and balance info
type KeeperInfo struct {
	Payee   string
	Active  bool
	Balance *big.Int
}

// UpkeepInfo keeper target info
type UpkeepInfo struct {
	Target                 string
	ExecuteGas             uint32
	CheckData              []byte
	Balance                *big.Int
	LastKeeper             string
	Admin                  string
	MaxValidBlocknumber    uint64
	LastPerformBlockNumber uint32
	AmountSpent            *big.Int
	Paused                 bool
	OffchainConfig         []byte
}

// LegacyEthereumKeeperRegistry represents keeper registry contract
type LegacyEthereumKeeperRegistry struct {
	client      blockchain.EVMClient
	version     ethereum.KeeperRegistryVersion
	registry1_1 *keeper_registry_wrapper1_1.KeeperRegistry
	registry1_2 *keeper_registry_wrapper1_2.KeeperRegistry
	registry1_3 *keeper_registry_wrapper1_3.KeeperRegistry
	registry2_0 *keeper_registry_wrapper2_0.KeeperRegistry
	registry2_1 *i_keeper_registry_master_wrapper_2_1.IKeeperRegistryMaster
	registry2_2 *i_automation_registry_master_wrapper_2_2.IAutomationRegistryMaster
	chainModule *i_chain_module.IChainModule
	address     *common.Address
	l           zerolog.Logger
}

func (v *LegacyEthereumKeeperRegistry) ReorgProtectionEnabled() bool {
	chainId := v.client.GetChainID().Uint64()
	// reorg protection is disabled in polygon zkEVM and Scroll bc currently there is no way to get the block hash onchain
	return v.version != ethereum.RegistryVersion_2_2 || (chainId != 1101 && chainId != 1442 && chainId != 2442 && chainId != 534352 && chainId != 534351)
}

func (v *LegacyEthereumKeeperRegistry) ChainModuleAddress() common.Address {
	if v.version == ethereum.RegistryVersion_2_2 {
		return v.chainModule.Address()
	}
	return common.Address{}
}

func (v *LegacyEthereumKeeperRegistry) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumKeeperRegistry) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(geth.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

func (rcs *KeeperRegistrySettings) Create22OnchainConfig(registrar string, registryOwnerAddress, chainModuleAddress common.Address, reorgProtectionEnabled bool) i_automation_registry_master_wrapper_2_2.AutomationRegistryBase22OnchainConfig {
	return i_automation_registry_master_wrapper_2_2.AutomationRegistryBase22OnchainConfig{
		PaymentPremiumPPB:      rcs.PaymentPremiumPPB,
		FlatFeeMicroLink:       rcs.FlatFeeMicroLINK,
		CheckGasLimit:          rcs.CheckGasLimit,
		StalenessSeconds:       rcs.StalenessSeconds,
		GasCeilingMultiplier:   rcs.GasCeilingMultiplier,
		MinUpkeepSpend:         rcs.MinUpkeepSpend,
		MaxPerformGas:          rcs.MaxPerformGas,
		MaxCheckDataSize:       rcs.MaxCheckDataSize,
		MaxPerformDataSize:     rcs.MaxPerformDataSize,
		MaxRevertDataSize:      rcs.MaxRevertDataSize,
		FallbackGasPrice:       rcs.FallbackGasPrice,
		FallbackLinkPrice:      rcs.FallbackLinkPrice,
		Transcoder:             common.Address{},
		Registrars:             []common.Address{common.HexToAddress(registrar)},
		UpkeepPrivilegeManager: registryOwnerAddress,
		ChainModule:            chainModuleAddress,
		ReorgProtectionEnabled: reorgProtectionEnabled,
	}
}

func (rcs *KeeperRegistrySettings) Create21OnchainConfig(registrar string, registryOwnerAddress common.Address) i_keeper_registry_master_wrapper_2_1.IAutomationV21PlusCommonOnchainConfigLegacy {
	return i_keeper_registry_master_wrapper_2_1.IAutomationV21PlusCommonOnchainConfigLegacy{
		PaymentPremiumPPB:      rcs.PaymentPremiumPPB,
		FlatFeeMicroLink:       rcs.FlatFeeMicroLINK,
		CheckGasLimit:          rcs.CheckGasLimit,
		StalenessSeconds:       rcs.StalenessSeconds,
		GasCeilingMultiplier:   rcs.GasCeilingMultiplier,
		MinUpkeepSpend:         rcs.MinUpkeepSpend,
		MaxPerformGas:          rcs.MaxPerformGas,
		MaxCheckDataSize:       rcs.MaxCheckDataSize,
		MaxPerformDataSize:     rcs.MaxPerformDataSize,
		MaxRevertDataSize:      rcs.MaxRevertDataSize,
		FallbackGasPrice:       rcs.FallbackGasPrice,
		FallbackLinkPrice:      rcs.FallbackLinkPrice,
		Transcoder:             common.Address{},
		Registrars:             []common.Address{common.HexToAddress(registrar)},
		UpkeepPrivilegeManager: registryOwnerAddress,
	}
}

func (rcs *KeeperRegistrySettings) Encode20OnchainConfig(registrar string) []byte {
	configType := goabi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address registrar)")
	onchainConfig, _ := goabi.Encode(map[string]interface{}{
		"paymentPremiumPPB":    rcs.PaymentPremiumPPB,
		"flatFeeMicroLink":     rcs.FlatFeeMicroLINK,
		"checkGasLimit":        rcs.CheckGasLimit,
		"stalenessSeconds":     rcs.StalenessSeconds,
		"gasCeilingMultiplier": rcs.GasCeilingMultiplier,
		"minUpkeepSpend":       rcs.MinUpkeepSpend,
		"maxPerformGas":        rcs.MaxPerformGas,
		"maxCheckDataSize":     rcs.MaxCheckDataSize,
		"maxPerformDataSize":   rcs.MaxPerformDataSize,
		"fallbackGasPrice":     rcs.FallbackGasPrice,
		"fallbackLinkPrice":    rcs.FallbackLinkPrice,
		"transcoder":           common.Address{},
		"registrar":            registrar,
	}, configType)
	return onchainConfig
}

func (v *LegacyEthereumKeeperRegistry) RegistryOwnerAddress() common.Address {
	callOpts := &bind.CallOpts{
		Pending: false,
	}

	//nolint: exhaustive
	switch v.version {
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
		return common.HexToAddress(v.client.GetDefaultWallet().Address())
	}

	return common.HexToAddress(v.client.GetDefaultWallet().Address())
}

func (v *LegacyEthereumKeeperRegistry) SetConfigTypeSafe(ocrConfig OCRv2Config) error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	switch v.version {
	case ethereum.RegistryVersion_2_1:
		tx, err := v.registry2_1.SetConfigTypeSafe(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.TypedOnchainConfig21,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		tx, err := v.registry2_2.SetConfigTypeSafe(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.TypedOnchainConfig22,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("SetConfigTypeSafe is not supported in keeper registry version %d", v.version)
	}
}

func (v *LegacyEthereumKeeperRegistry) SetConfig(config KeeperRegistrySettings, ocrConfig OCRv2Config) error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	callOpts := bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: nil,
	}
	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err := v.registry1_1.SetConfig(
			txOpts,
			config.PaymentPremiumPPB,
			config.FlatFeeMicroLINK,
			config.BlockCountPerTurn,
			config.CheckGasLimit,
			config.StalenessSeconds,
			config.GasCeilingMultiplier,
			config.FallbackGasPrice,
			config.FallbackLinkPrice,
		)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_1_2:
		state, err := v.registry1_2.GetState(&callOpts)
		if err != nil {
			return err
		}

		tx, err := v.registry1_2.SetConfig(txOpts, keeper_registry_wrapper1_2.Config{
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
		})
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_1_3:
		state, err := v.registry1_3.GetState(&callOpts)
		if err != nil {
			return err
		}

		tx, err := v.registry1_3.SetConfig(txOpts, keeper_registry_wrapper1_3.Config{
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
		})
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_0:
		tx, err := v.registry2_0.SetConfig(txOpts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.OnchainConfig,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2:
		return fmt.Errorf("registry version 2.1 and 2.2 must use setConfigTypeSafe function")
	default:
		return fmt.Errorf("keeper registry version %d is not supported", v.version)
	}
}

// Pause pauses the registry.
func (v *LegacyEthereumKeeperRegistry) Pause() error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	var tx *types.Transaction

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err = v.registry1_1.Pause(txOpts)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_1_2:
		tx, err = v.registry1_2.Pause(txOpts)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_1_3:
		tx, err = v.registry1_3.Pause(txOpts)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_0:
		tx, err = v.registry2_0.Pause(txOpts)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_1:
		tx, err = v.registry2_1.Pause(txOpts)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		tx, err = v.registry2_2.Pause(txOpts)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	}

	return fmt.Errorf("keeper registry version %d is not supported", v.version)
}

// Migrate performs a migration of the given upkeep ids to the specific destination passed as parameter.
func (v *LegacyEthereumKeeperRegistry) Migrate(upkeepIDs []*big.Int, destinationAddress common.Address) error {
	if v.version != ethereum.RegistryVersion_1_2 {
		return fmt.Errorf("migration of upkeeps is only available for version 1.2 of the registries")
	}

	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := v.registry1_2.MigrateUpkeeps(txOpts, upkeepIDs, destinationAddress)
	if err != nil {
		return err
	}

	return v.client.ProcessTransaction(tx)
}

// SetMigrationPermissions sets the permissions of another registry to allow migrations between the two.
func (v *LegacyEthereumKeeperRegistry) SetMigrationPermissions(peerAddress common.Address, permission uint8) error {
	if v.version != ethereum.RegistryVersion_1_2 {
		return fmt.Errorf("migration of upkeeps is only available for version 1.2 of the registries")
	}

	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := v.registry1_2.SetPeerRegistryMigrationPermission(txOpts, peerAddress, permission)
	if err != nil {
		return err
	}

	return v.client.ProcessTransaction(tx)
}

func (v *LegacyEthereumKeeperRegistry) SetRegistrar(registrarAddr string) error {
	if v.version == ethereum.RegistryVersion_2_0 {
		// we short circuit and exit, so we don't create a new txs messing up the nonce before exiting
		return fmt.Errorf("please use set config")
	}

	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	callOpts := bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: nil,
	}

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err := v.registry1_1.SetRegistrar(txOpts, common.HexToAddress(registrarAddr))
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_1_2:
		state, err := v.registry1_2.GetState(&callOpts)
		if err != nil {
			return err
		}
		newConfig := state.Config
		newConfig.Registrar = common.HexToAddress(registrarAddr)
		tx, err := v.registry1_2.SetConfig(txOpts, newConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_1_3:
		state, err := v.registry1_3.GetState(&callOpts)
		if err != nil {
			return err
		}
		newConfig := state.Config
		newConfig.Registrar = common.HexToAddress(registrarAddr)
		tx, err := v.registry1_3.SetConfig(txOpts, newConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("keeper registry version %d is not supported", v.version)
	}
}

// AddUpkeepFunds adds link for particular upkeep id
func (v *LegacyEthereumKeeperRegistry) AddUpkeepFundsFromKey(_ *big.Int, _ *big.Int, _ int) error {
	panic("this method is only supported by contracts using Seth client")
}

// AddUpkeepFunds adds link for particular upkeep id
func (v *LegacyEthereumKeeperRegistry) AddUpkeepFunds(id *big.Int, amount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	var tx *types.Transaction

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err = v.registry1_1.AddFunds(opts, id, amount)
	case ethereum.RegistryVersion_1_2:
		tx, err = v.registry1_2.AddFunds(opts, id, amount)
	case ethereum.RegistryVersion_1_3:
		tx, err = v.registry1_3.AddFunds(opts, id, amount)
	case ethereum.RegistryVersion_2_0:
		tx, err = v.registry2_0.AddFunds(opts, id, amount)
	case ethereum.RegistryVersion_2_1:
		tx, err = v.registry2_1.AddFunds(opts, id, amount)
	case ethereum.RegistryVersion_2_2:
		tx, err = v.registry2_2.AddFunds(opts, id, amount)
	}

	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// GetUpkeepInfo gets upkeep info
func (v *LegacyEthereumKeeperRegistry) GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
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
	}

	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

func (v *LegacyEthereumKeeperRegistry) getUpkeepInfo22(opts *bind.CallOpts, id *big.Int) (*UpkeepInfo, error) {
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

func (v *LegacyEthereumKeeperRegistry) GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
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
	case ethereum.RegistryVersion_2_0, ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2:
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

func (v *LegacyEthereumKeeperRegistry) SetKeepers(keepers []string, payees []string, ocrConfig OCRv2Config) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	keepersAddresses := make([]common.Address, 0)
	for _, k := range keepers {
		keepersAddresses = append(keepersAddresses, common.HexToAddress(k))
	}
	payeesAddresses := make([]common.Address, 0)
	for _, p := range payees {
		payeesAddresses = append(payeesAddresses, common.HexToAddress(p))
	}
	var tx *types.Transaction

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err = v.registry1_1.SetKeepers(opts, keepersAddresses, payeesAddresses)
	case ethereum.RegistryVersion_1_2:
		tx, err = v.registry1_2.SetKeepers(opts, keepersAddresses, payeesAddresses)
	case ethereum.RegistryVersion_1_3:
		tx, err = v.registry1_3.SetKeepers(opts, keepersAddresses, payeesAddresses)
	case ethereum.RegistryVersion_2_0:
		tx, err = v.registry2_0.SetConfig(opts,
			ocrConfig.Signers,
			ocrConfig.Transmitters,
			ocrConfig.F,
			ocrConfig.OnchainConfig,
			ocrConfig.OffchainConfigVersion,
			ocrConfig.OffchainConfig,
		)
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2:
		return fmt.Errorf("not supported")
	}

	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// RegisterUpkeep registers contract to perform upkeep
func (v *LegacyEthereumKeeperRegistry) RegisterUpkeep(target string, gasLimit uint32, admin string, checkData []byte) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	var tx *types.Transaction

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err = v.registry1_1.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
		)
	case ethereum.RegistryVersion_1_2:
		tx, err = v.registry1_2.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
		)
	case ethereum.RegistryVersion_1_3:
		tx, err = v.registry1_3.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
		)
	case ethereum.RegistryVersion_2_0:
		tx, err = v.registry2_0.RegisterUpkeep(
			opts,
			common.HexToAddress(target),
			gasLimit,
			common.HexToAddress(admin),
			checkData,
			nil, //offchain config
		)
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2:
		return fmt.Errorf("not supported")
	}

	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// CancelUpkeep cancels the given upkeep ID
func (v *LegacyEthereumKeeperRegistry) CancelUpkeep(id *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	var tx *types.Transaction

	switch v.version {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		tx, err = v.registry1_1.CancelUpkeep(opts, id)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_1_2:
		tx, err = v.registry1_2.CancelUpkeep(opts, id)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_1_3:
		tx, err = v.registry1_3.CancelUpkeep(opts, id)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_2_0:
		tx, err = v.registry2_0.CancelUpkeep(opts, id)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_2_1:
		tx, err = v.registry2_1.CancelUpkeep(opts, id)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_2_2:
		tx, err = v.registry2_2.CancelUpkeep(opts, id)
		if err != nil {
			return err
		}
	}

	v.l.Info().
		Str("Upkeep ID", strconv.FormatInt(id.Int64(), 10)).
		Str("From", v.client.GetDefaultWallet().Address()).
		Str("TX Hash", tx.Hash().String()).
		Msg("Cancel Upkeep tx")
	return v.client.ProcessTransaction(tx)
}

// SetUpkeepGasLimit sets the perform gas limit for a given upkeep ID
func (v *LegacyEthereumKeeperRegistry) SetUpkeepGasLimit(id *big.Int, gas uint32) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	var tx *types.Transaction

	switch v.version {
	case ethereum.RegistryVersion_1_2:
		tx, err = v.registry1_2.SetUpkeepGasLimit(opts, id, gas)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_1_3:
		tx, err = v.registry1_3.SetUpkeepGasLimit(opts, id, gas)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_2_0:
		tx, err = v.registry2_0.SetUpkeepGasLimit(opts, id, gas)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_2_1:
		tx, err = v.registry2_1.SetUpkeepGasLimit(opts, id, gas)
		if err != nil {
			return err
		}
	case ethereum.RegistryVersion_2_2:
		tx, err = v.registry2_2.SetUpkeepGasLimit(opts, id, gas)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("keeper registry version %d is not supported for SetUpkeepGasLimit", v.version)
	}
	return v.client.ProcessTransaction(tx)
}

// GetKeeperList get list of all registered keeper addresses
func (v *LegacyEthereumKeeperRegistry) GetKeeperList(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
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
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2:
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
func (v *LegacyEthereumKeeperRegistry) UpdateCheckData(id *big.Int, newCheckData []byte) error {

	switch v.version {
	case ethereum.RegistryVersion_1_3:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry1_3.UpdateCheckData(opts, id, newCheckData)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_0:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_0.UpdateCheckData(opts, id, newCheckData)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_1:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_1.SetUpkeepCheckData(opts, id, newCheckData)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_2.SetUpkeepCheckData(opts, id, newCheckData)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("UpdateCheckData is not supported by keeper registry version %d", v.version)
	}
}

// SetUpkeepTriggerConfig updates the trigger config of an upkeep (only for version 2.1)
func (v *LegacyEthereumKeeperRegistry) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) error {

	switch v.version {
	case ethereum.RegistryVersion_2_1:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_1.SetUpkeepTriggerConfig(opts, id, triggerConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_2.SetUpkeepTriggerConfig(opts, id, triggerConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("SetUpkeepTriggerConfig is not supported by keeper registry version %d", v.version)
	}
}

// SetUpkeepPrivilegeConfig sets the privilege config of an upkeep (only for version 2.1)
func (v *LegacyEthereumKeeperRegistry) SetUpkeepPrivilegeConfig(id *big.Int, privilegeConfig []byte) error {

	switch v.version {
	case ethereum.RegistryVersion_2_1:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_1.SetUpkeepPrivilegeConfig(opts, id, privilegeConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_2.SetUpkeepPrivilegeConfig(opts, id, privilegeConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("SetUpkeepPrivilegeConfig is not supported by keeper registry version %d", v.version)
	}
}

// PauseUpkeep stops an upkeep from an upkeep
func (v *LegacyEthereumKeeperRegistry) PauseUpkeep(id *big.Int) error {
	switch v.version {
	case ethereum.RegistryVersion_1_3:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry1_3.PauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_0:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_0.PauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_1:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_1.PauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_2.PauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("PauseUpkeep is not supported by keeper registry version %d", v.version)
	}
}

// UnpauseUpkeep get list of all registered keeper addresses
func (v *LegacyEthereumKeeperRegistry) UnpauseUpkeep(id *big.Int) error {
	switch v.version {
	case ethereum.RegistryVersion_1_3:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry1_3.UnpauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_0:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_0.UnpauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_1:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_1.UnpauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_2.UnpauseUpkeep(opts, id)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("UnpauseUpkeep is not supported by keeper registry version %d", v.version)
	}
}

func (v *LegacyEthereumKeeperRegistry) SetUpkeepOffchainConfig(id *big.Int, offchainConfig []byte) error {
	switch v.version {
	case ethereum.RegistryVersion_2_0:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_0.SetUpkeepOffchainConfig(opts, id, offchainConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_1:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_1.SetUpkeepOffchainConfig(opts, id, offchainConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	case ethereum.RegistryVersion_2_2:
		opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
		if err != nil {
			return err
		}

		tx, err := v.registry2_2.SetUpkeepOffchainConfig(opts, id, offchainConfig)
		if err != nil {
			return err
		}
		return v.client.ProcessTransaction(tx)
	default:
		return fmt.Errorf("SetUpkeepOffchainConfig is not supported by keeper registry version %d", v.version)
	}
}

// Parses upkeep performed log
func (v *LegacyEthereumKeeperRegistry) ParseUpkeepPerformedLog(log *types.Log) (*UpkeepPerformedLog, error) {
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
	}
	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

// ParseStaleUpkeepReportLog Parses Stale upkeep report log
func (v *LegacyEthereumKeeperRegistry) ParseStaleUpkeepReportLog(log *types.Log) (*StaleUpkeepReportLog, error) {
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
	}
	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

// Parses the upkeep ID from an 'UpkeepRegistered' log, returns error on any other log
func (v *LegacyEthereumKeeperRegistry) ParseUpkeepIdFromRegisteredLog(log *types.Log) (*big.Int, error) {
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
	}

	return nil, fmt.Errorf("keeper registry version %d is not supported", v.version)
}

// KeeperConsumerRoundConfirmer is a header subscription that awaits for a round of upkeeps
type KeeperConsumerRoundConfirmer struct {
	instance     KeeperConsumer
	upkeepsValue int
	doneChan     chan struct{}
	context      context.Context
	cancel       context.CancelFunc
	l            zerolog.Logger
}

// NewKeeperConsumerRoundConfirmer provides a new instance of a KeeperConsumerRoundConfirmer
func NewKeeperConsumerRoundConfirmer(
	contract KeeperConsumer,
	counterValue int,
	timeout time.Duration,
	logger zerolog.Logger,
) *KeeperConsumerRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &KeeperConsumerRoundConfirmer{
		instance:     contract,
		upkeepsValue: counterValue,
		doneChan:     make(chan struct{}),
		context:      ctx,
		cancel:       ctxCancel,
		l:            logger,
	}
}

// ReceiveHeader will query the latest Keeper round and check to see whether the round has confirmed
func (o *KeeperConsumerRoundConfirmer) ReceiveHeader(_ blockchain.NodeHeader) error {
	upkeeps, err := o.instance.Counter(context.Background())
	if err != nil {
		return err
	}
	l := o.l.Info().
		Str("Contract Address", o.instance.Address()).
		Int64("Upkeeps", upkeeps.Int64()).
		Int("Required upkeeps", o.upkeepsValue)
	if upkeeps.Int64() == int64(o.upkeepsValue) {
		l.Msg("Upkeep completed")
		o.doneChan <- struct{}{}
	} else {
		l.Msg("Waiting for upkeep round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *KeeperConsumerRoundConfirmer) Wait() error {
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for upkeeps to confirm: %d", o.upkeepsValue)
		}
	}
}

// KeeperConsumerPerformanceRoundConfirmer is a header subscription that awaits for a round of upkeeps
type KeeperConsumerPerformanceRoundConfirmer struct {
	instance KeeperConsumerPerformance
	doneChan chan bool
	context  context.Context
	cancel   context.CancelFunc

	lastBlockNum                uint64  // Records the number of the last block that came in
	blockCadence                int64   // How many blocks before an upkeep should happen
	blockRange                  int64   // How many blocks to watch upkeeps for
	blocksSinceSubscription     int64   // How many blocks have passed since subscribing
	expectedUpkeepCount         int64   // The count of upkeeps expected next iteration
	blocksSinceSuccessfulUpkeep int64   // How many blocks have come in since the last successful upkeep
	allMissedUpkeeps            []int64 // Tracks the amount of blocks missed in each missed upkeep
	totalSuccessfulUpkeeps      int64

	metricsReporter *testreporters.KeeperBlockTimeTestReporter // Testreporter to track results
	complete        bool
	l               zerolog.Logger
}

// NewKeeperConsumerPerformanceRoundConfirmer provides a new instance of a KeeperConsumerPerformanceRoundConfirmer
// Used to track and log performance test results for keepers
func NewKeeperConsumerPerformanceRoundConfirmer(
	contract KeeperConsumerPerformance,
	expectedBlockCadence int64, // Expected to upkeep every 5/10/20 blocks, for example
	blockRange int64,
	metricsReporter *testreporters.KeeperBlockTimeTestReporter,
	logger zerolog.Logger,
) *KeeperConsumerPerformanceRoundConfirmer {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &KeeperConsumerPerformanceRoundConfirmer{
		instance:                    contract,
		doneChan:                    make(chan bool),
		context:                     ctx,
		cancel:                      cancelFunc,
		blockCadence:                expectedBlockCadence,
		blockRange:                  blockRange,
		blocksSinceSubscription:     0,
		blocksSinceSuccessfulUpkeep: 0,
		expectedUpkeepCount:         1,
		allMissedUpkeeps:            []int64{},
		totalSuccessfulUpkeeps:      0,
		metricsReporter:             metricsReporter,
		complete:                    false,
		lastBlockNum:                0,
		l:                           logger,
	}
}

// ReceiveHeader will query the latest Keeper round and check to see whether the round has confirmed
func (o *KeeperConsumerPerformanceRoundConfirmer) ReceiveHeader(receivedHeader blockchain.NodeHeader) error {
	if receivedHeader.Number.Uint64() <= o.lastBlockNum { // Uncle / reorg we won't count
		return nil
	}
	o.lastBlockNum = receivedHeader.Number.Uint64()
	// Increment block counters
	o.blocksSinceSubscription++
	o.blocksSinceSuccessfulUpkeep++
	upkeepCount, err := o.instance.GetUpkeepCount(context.Background())
	if err != nil {
		return err
	}

	isEligible, err := o.instance.CheckEligible(context.Background())
	if err != nil {
		return err
	}
	if isEligible {
		o.l.Trace().
			Str("Contract Address", o.instance.Address()).
			Int64("Upkeeps Performed", upkeepCount.Int64()).
			Msg("Upkeep Now Eligible")
	}
	if upkeepCount.Int64() >= o.expectedUpkeepCount { // Upkeep was successful
		if o.blocksSinceSuccessfulUpkeep < o.blockCadence { // If there's an early upkeep, that's weird
			o.l.Error().
				Str("Contract Address", o.instance.Address()).
				Int64("Upkeeps Performed", upkeepCount.Int64()).
				Int64("Expected Cadence", o.blockCadence).
				Int64("Actual Cadence", o.blocksSinceSuccessfulUpkeep).
				Err(errors.New("found an early Upkeep"))
			return fmt.Errorf("found an early Upkeep on contract %s", o.instance.Address())
		} else if o.blocksSinceSuccessfulUpkeep == o.blockCadence { // Perfectly timed upkeep
			o.l.Info().
				Str("Contract Address", o.instance.Address()).
				Int64("Upkeeps Performed", upkeepCount.Int64()).
				Int64("Expected Cadence", o.blockCadence).
				Int64("Actual Cadence", o.blocksSinceSuccessfulUpkeep).
				Msg("Successful Upkeep on Expected Cadence")
			o.totalSuccessfulUpkeeps++
		} else { // Late upkeep
			o.l.Warn().
				Str("Contract Address", o.instance.Address()).
				Int64("Upkeeps Performed", upkeepCount.Int64()).
				Int64("Expected Cadence", o.blockCadence).
				Int64("Actual Cadence", o.blocksSinceSuccessfulUpkeep).
				Msg("Upkeep Completed Late")
			o.allMissedUpkeeps = append(o.allMissedUpkeeps, o.blocksSinceSuccessfulUpkeep-o.blockCadence)
		}
		// Update upkeep tracking values
		o.blocksSinceSuccessfulUpkeep = 0
		o.expectedUpkeepCount++
	}

	if o.blocksSinceSubscription > o.blockRange {
		if o.blocksSinceSuccessfulUpkeep > o.blockCadence {
			o.l.Warn().
				Str("Contract Address", o.instance.Address()).
				Int64("Upkeeps Performed", upkeepCount.Int64()).
				Int64("Expected Cadence", o.blockCadence).
				Int64("Expected Upkeep Count", o.expectedUpkeepCount).
				Int64("Blocks Waiting", o.blocksSinceSuccessfulUpkeep).
				Int64("Total Blocks Watched", o.blocksSinceSubscription).
				Msg("Finished Watching for Upkeeps While Waiting on a Late Upkeep")
			o.allMissedUpkeeps = append(o.allMissedUpkeeps, o.blocksSinceSuccessfulUpkeep-o.blockCadence)
		} else {
			o.l.Info().
				Str("Contract Address", o.instance.Address()).
				Int64("Upkeeps Performed", upkeepCount.Int64()).
				Int64("Total Blocks Watched", o.blocksSinceSubscription).
				Msg("Finished Watching for Upkeeps")
		}
		o.doneChan <- true
		o.complete = true
		return nil
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *KeeperConsumerPerformanceRoundConfirmer) Wait() error {
	defer func() { o.complete = true }()
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			o.logDetails()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for expected upkeep count to confirm: %d", o.expectedUpkeepCount)
		}
	}
}

func (o *KeeperConsumerPerformanceRoundConfirmer) Complete() bool {
	return o.complete
}

func (o *KeeperConsumerPerformanceRoundConfirmer) logDetails() {
	report := testreporters.KeeperBlockTimeTestReport{
		ContractAddress:        o.instance.Address(),
		TotalExpectedUpkeeps:   o.blockRange / o.blockCadence,
		TotalSuccessfulUpkeeps: o.totalSuccessfulUpkeeps,
		AllMissedUpkeeps:       o.allMissedUpkeeps,
	}
	o.metricsReporter.ReportMutex.Lock()
	o.metricsReporter.Reports = append(o.metricsReporter.Reports, report)
	defer o.metricsReporter.ReportMutex.Unlock()
}

// LegacyKeeperConsumerBenchmarkRoundConfirmer is a header subscription that awaits for a round of upkeeps
type LegacyKeeperConsumerBenchmarkRoundConfirmer struct {
	instance AutomationConsumerBenchmark
	registry KeeperRegistry
	upkeepID *big.Int
	doneChan chan bool
	context  context.Context
	cancel   context.CancelFunc

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

// NewLegacyKeeperConsumerBenchmarkRoundConfirmer provides a new instance of a LegacyKeeperConsumerBenchmarkRoundConfirmer
// Used to track and log benchmark test results for keepers
func NewLegacyKeeperConsumerBenchmarkRoundConfirmer(
	contract AutomationConsumerBenchmark,
	registry KeeperRegistry,
	upkeepID *big.Int,
	blockRange int64,
	upkeepSLA int64,
	metricsReporter *testreporters.KeeperBenchmarkTestReporter,
	upkeepIndex int64,
	firstEligibleBuffer int64,
	logger zerolog.Logger,
) *LegacyKeeperConsumerBenchmarkRoundConfirmer {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &LegacyKeeperConsumerBenchmarkRoundConfirmer{
		instance:                contract,
		registry:                registry,
		upkeepID:                upkeepID,
		doneChan:                make(chan bool),
		context:                 ctx,
		cancel:                  cancelFunc,
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

// ReceiveHeader will query the latest Keeper round and check to see whether the round has confirmed
func (o *LegacyKeeperConsumerBenchmarkRoundConfirmer) ReceiveHeader(receivedHeader blockchain.NodeHeader) error {
	if receivedHeader.Number.Uint64() <= o.lastBlockNum { // Uncle / reorg we won't count
		return nil
	}
	if o.firstBlockNum == 0 {
		o.firstBlockNum = receivedHeader.Number.Uint64()
	}
	o.lastBlockNum = receivedHeader.Number.Uint64()
	// Increment block counters
	o.blocksSinceSubscription++

	upkeepCount, err := o.instance.GetUpkeepCount(context.Background(), big.NewInt(o.upkeepIndex))
	if err != nil {
		return err
	}

	if upkeepCount.Int64() > o.upkeepCount { // A new upkeep was done
		if upkeepCount.Int64() != o.upkeepCount+1 {
			return errors.New("upkeep count increased by more than 1 in a single block")
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
		return err
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

		o.doneChan <- true
		o.complete = true
		return nil
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *LegacyKeeperConsumerBenchmarkRoundConfirmer) Wait() error {
	defer func() { o.complete = true }()
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			o.logDetails()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for expected number of blocks: %d", o.blockRange)
		}
	}
}

func (o *LegacyKeeperConsumerBenchmarkRoundConfirmer) Complete() bool {
	return o.complete
}

func (o *LegacyKeeperConsumerBenchmarkRoundConfirmer) logDetails() {
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

// LegacyEthereumUpkeepCounter represents keeper consumer (upkeep) counter contract
type LegacyEthereumUpkeepCounter struct {
	client   blockchain.EVMClient
	consumer *upkeep_counter_wrapper.UpkeepCounter
	address  *common.Address
}

func (v *LegacyEthereumUpkeepCounter) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumUpkeepCounter) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(geth.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}
func (v *LegacyEthereumUpkeepCounter) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

func (v *LegacyEthereumUpkeepCounter) SetSpread(testRange *big.Int, interval *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.SetSpread(opts, testRange, interval)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// Just pass for non-logtrigger
func (v *LegacyEthereumUpkeepCounter) Start() error {
	return nil
}

// LegacyEthereumUpkeepPerformCounterRestrictive represents keeper consumer (upkeep) counter contract
type LegacyEthereumUpkeepPerformCounterRestrictive struct {
	client   blockchain.EVMClient
	consumer *upkeep_perform_counter_restrictive_wrapper.UpkeepPerformCounterRestrictive
	address  *common.Address
}

func (v *LegacyEthereumUpkeepPerformCounterRestrictive) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumUpkeepPerformCounterRestrictive) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(geth.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}
func (v *LegacyEthereumUpkeepPerformCounterRestrictive) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	count, err := v.consumer.GetCountPerforms(opts)
	return count, err
}

func (v *LegacyEthereumUpkeepPerformCounterRestrictive) SetSpread(testRange *big.Int, interval *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.SetSpread(opts, testRange, interval)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// EthereumKeeperConsumer represents keeper consumer (upkeep) contract
type EthereumKeeperConsumer struct {
	client   blockchain.EVMClient
	consumer *keeper_consumer_wrapper.KeeperConsumer
	address  *common.Address
}

// Just pass for non-logtrigger
func (v *EthereumKeeperConsumer) Start() error {
	return nil
}

func (v *EthereumKeeperConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

type LegacyEthereumAutomationStreamsLookupUpkeepConsumer struct {
	client   blockchain.EVMClient
	consumer *streams_lookup_upkeep_wrapper.StreamsLookupUpkeep
	address  *common.Address
}

func (v *LegacyEthereumAutomationStreamsLookupUpkeepConsumer) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumAutomationStreamsLookupUpkeepConsumer) Start() error {
	// For this consumer upkeep, we use this Start() function to set ParamKeys so as to run mercury v0.2
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	// The default values of ParamKeys are "feedIDs" and "timestamp" which are for v0.3
	tx, err := v.consumer.SetParamKeys(txOpts, "feedIdHex", "blockNumber")
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *LegacyEthereumAutomationStreamsLookupUpkeepConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

type LegacyEthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer struct {
	client   blockchain.EVMClient
	consumer *log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookup
	address  *common.Address
}

func (v *LegacyEthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer) Address() string {
	return v.address.Hex()
}

// Kick off the log trigger event. The contract uses Mercury v0.2 so no need to set ParamKeys
func (v *LegacyEthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer) Start() error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := v.consumer.Start(txOpts)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *LegacyEthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

type LegacyEthereumAutomationLogCounterConsumer struct {
	client   blockchain.EVMClient
	consumer *log_upkeep_counter_wrapper.LogUpkeepCounter
	address  *common.Address
}

func (v *LegacyEthereumAutomationLogCounterConsumer) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumAutomationLogCounterConsumer) Start() error {
	txOpts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := v.consumer.Start(txOpts)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *LegacyEthereumAutomationLogCounterConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

type LegacyEthereumAutomationSimpleLogCounterConsumer struct {
	client   blockchain.EVMClient
	consumer *simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounter
	address  *common.Address
}

func (v *LegacyEthereumAutomationSimpleLogCounterConsumer) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumAutomationSimpleLogCounterConsumer) Start() error {
	return nil
}

func (v *LegacyEthereumAutomationSimpleLogCounterConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

// LegacyEthereumKeeperConsumerPerformance represents a more complicated keeper consumer contract, one intended only for
// performance tests.
type LegacyEthereumKeeperConsumerPerformance struct {
	client   blockchain.EVMClient
	consumer *keeper_consumer_performance_wrapper.KeeperConsumerPerformance
	address  *common.Address
}

func (v *LegacyEthereumKeeperConsumerPerformance) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumKeeperConsumerPerformance) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(geth.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

func (v *LegacyEthereumKeeperConsumerPerformance) CheckEligible(ctx context.Context) (bool, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	eligible, err := v.consumer.CheckEligible(opts)
	return eligible, err
}

func (v *LegacyEthereumKeeperConsumerPerformance) GetUpkeepCount(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	eligible, err := v.consumer.GetCountPerforms(opts)
	return eligible, err
}

func (v *LegacyEthereumKeeperConsumerPerformance) SetCheckGasToBurn(_ context.Context, gas *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.SetCheckGasToBurn(opts, gas)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *LegacyEthereumKeeperConsumerPerformance) SetPerformGasToBurn(_ context.Context, gas *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.SetPerformGasToBurn(opts, gas)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// LegacyEthereumKeeperPerformDataCheckerConsumer represents keeper perform data checker contract
type LegacyEthereumKeeperPerformDataCheckerConsumer struct {
	client             blockchain.EVMClient
	performDataChecker *perform_data_checker_wrapper.PerformDataChecker
	address            *common.Address
}

func (v *LegacyEthereumKeeperPerformDataCheckerConsumer) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumKeeperPerformDataCheckerConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	cnt, err := v.performDataChecker.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

func (v *LegacyEthereumKeeperPerformDataCheckerConsumer) SetExpectedData(_ context.Context, expectedData []byte) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.performDataChecker.SetExpectedData(opts, expectedData)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// LegacyEthereumAutomationConsumerBenchmark represents a more complicated keeper consumer contract, one intended only for
// Benchmark tests.
type LegacyEthereumAutomationConsumerBenchmark struct {
	client   blockchain.EVMClient
	consumer *automation_consumer_benchmark.AutomationConsumerBenchmark
	address  *common.Address
}

func (v *LegacyEthereumAutomationConsumerBenchmark) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumAutomationConsumerBenchmark) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(geth.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

func (v *LegacyEthereumAutomationConsumerBenchmark) CheckEligible(ctx context.Context, id *big.Int, _range *big.Int, firstEligibleBuffer *big.Int) (bool, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	eligible, err := v.consumer.CheckEligible(opts, id, _range, firstEligibleBuffer)
	return eligible, err
}

func (v *LegacyEthereumAutomationConsumerBenchmark) GetUpkeepCount(ctx context.Context, id *big.Int) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	eligible, err := v.consumer.GetCountPerforms(opts, id)
	return eligible, err
}

// LegacyEthereumKeeperRegistrar corresponds to the registrar which is used to send requests to the registry when
// registering new upkeeps.
type LegacyEthereumKeeperRegistrar struct {
	client      blockchain.EVMClient
	registrar   *keeper_registrar_wrapper1_2.KeeperRegistrar
	registrar20 *keeper_registrar_wrapper2_0.KeeperRegistrar
	registrar21 *registrar21.AutomationRegistrar
	address     *common.Address
}

func (v *LegacyEthereumKeeperRegistrar) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumKeeperRegistrar) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(geth.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

// EncodeRegisterRequest encodes register request to call it through link token TransferAndCall
func (v *LegacyEthereumKeeperRegistrar) EncodeRegisterRequest(name string, email []byte, upkeepAddr string, gasLimit uint32, adminAddr string, checkData []byte, amount *big.Int, source uint8, senderAddr string, isLogTrigger bool, isMercury bool) ([]byte, error) {
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

			logTriggerConfigStruct := ac.IAutomationV21PlusCommonLogTriggerConfig{
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

// LegacyEthereumUpkeepTranscoder represents the transcoder which is used to perform migrations
// of upkeeps from one registry to another.
type LegacyEthereumUpkeepTranscoder struct {
	client     blockchain.EVMClient
	transcoder *upkeep_transcoder.UpkeepTranscoder
	address    *common.Address
}

func (v *LegacyEthereumUpkeepTranscoder) Address() string {
	return v.address.Hex()
}
