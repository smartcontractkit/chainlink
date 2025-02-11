package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	goabi "github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"

	cltypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

var compatibleUtils = cltypes.MustGetABI(ac.AutomationCompatibleUtilsABI)
var registrarABI = cltypes.MustGetABI(registrar21.AutomationRegistrarABI)

type KeeperRegistrar interface {
	Address() string

	EncodeRegisterRequest(name string, email []byte, upkeepAddr string, gasLimit uint32, adminAddr string, checkData []byte, amount *big.Int, source uint8, senderAddr string, isLogTrigger bool, isMercury bool, linkTokenAddr string) ([]byte, error)

	Fund(ethAmount *big.Float) error

	RegisterUpkeepFromKey(keyNum int, name string, email []byte, upkeepAddr string, gasLimit uint32, adminAddr string, checkData []byte, amount *big.Int, wethTokenAddr string, isLogTrigger bool, isMercury bool) (*types.Transaction, error)
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
	RegistryVersion   ethereum.KeeperRegistryVersion
	LinkAddr          string
	ETHFeedAddr       string
	GasFeedAddr       string
	TranscoderAddr    string
	RegistrarAddr     string
	Settings          KeeperRegistrySettings
	LinkUSDFeedAddr   string
	NativeUSDFeedAddr string
	WrappedNativeAddr string
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
	FallbackNativePrice  *big.Int // Native price used if the Native price feed is stale
	MaxCheckDataSize     uint32
	MaxPerformDataSize   uint32
	MaxRevertDataSize    uint32
	RegistryVersion      ethereum.KeeperRegistryVersion
}

func (rcs *KeeperRegistrySettings) Create23OnchainConfig(registrar string, registryOwnerAddress, chainModuleAddress common.Address, reorgProtectionEnabled bool) i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23OnchainConfig {
	return i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23OnchainConfig{
		CheckGasLimit:          rcs.CheckGasLimit,
		StalenessSeconds:       rcs.StalenessSeconds,
		GasCeilingMultiplier:   rcs.GasCeilingMultiplier,
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
		FinanceAdmin:           registryOwnerAddress,
		FallbackNativePrice:    rcs.FallbackNativePrice,
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

// KeeperRegistrarSettings represents settings for registrar contract
type KeeperRegistrarSettings struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint16
	RegistryAddr          string
	MinLinkJuels          *big.Int
	WETHTokenAddr         string
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
