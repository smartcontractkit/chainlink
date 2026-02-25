package contracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	goabi "github.com/umbracle/ethgo/abi"

	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_compatible_utils"
	registrar21 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_consumer_performance_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/perform_data_checker_wrapper"
	cltypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
)

// AbigenLog is an interface for abigen generated log topics
type AbigenLog interface {
	Topic() common.Hash
}

type KeeperRegistryVersion int32

//nolint:revive // we want to use underscores
const (
	RegistryVersion_1_0 KeeperRegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
	RegistryVersion_1_3
	RegistryVersion_2_0
	RegistryVersion_2_1
	RegistryVersion_2_2
	RegistryVersion_2_3
)

func (k KeeperRegistryVersion) String() string {
	switch k {
	case RegistryVersion_1_0:
		return "1.0"
	case RegistryVersion_1_1:
		return "1.1"
	case RegistryVersion_1_2:
		return "1.2"
	case RegistryVersion_1_3:
		return "1.3"
	case RegistryVersion_2_0:
		return "2.0"
	case RegistryVersion_2_1:
		return "2.1"
	case RegistryVersion_2_2:
		return "2.2"
	case RegistryVersion_2_3:
		return "2.3"
	default:
		return "unknown"
	}
}

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
	ParseUpkeepIDFromRegisteredLog(log *types.Log) (*big.Int, error)
	Pause() error
	Unpause() error
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

type UpkeepPerformedLog struct {
	ID      *big.Int
	Success bool
	From    common.Address
}

type StaleUpkeepReportLog struct {
	ID *big.Int
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

// KeeperRegistryOpts opts to deploy keeper registry version
type KeeperRegistryOpts struct {
	RegistryVersion   KeeperRegistryVersion
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
	RegistryVersion      KeeperRegistryVersion
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
	onchainConfig, _ := goabi.Encode(map[string]any{
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

	instance, err := keeper_consumer_performance_wrapper.NewKeeperConsumerPerformance(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperConsumerPerformance{}, fmt.Errorf("failed to instantiate KeeperConsumerPerformance instance: %w", err)
	}

	return &EthereumKeeperConsumerPerformance{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}

type KeeperPerformDataChecker interface {
	Address() string
	Counter(ctx context.Context) (*big.Int, error)
	SetExpectedData(ctx context.Context, expectedData []byte) error
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

	instance, err := perform_data_checker_wrapper.NewPerformDataChecker(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperPerformDataCheckerConsumer{}, fmt.Errorf("failed to instantiate PerformDataChecker instance: %w", err)
	}

	return &EthereumKeeperPerformDataCheckerConsumer{
		client:             client,
		performDataChecker: instance,
		address:            &data.Address,
	}, nil
}
