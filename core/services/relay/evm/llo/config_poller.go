package llo

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var _ ocr2types.ContractConfigTracker = (*lloConfigPoller)(nil)

type LLOLogPoller interface {
	RegisterFilter(ctx context.Context, filter logpoller.Filter) error
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
	LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error)
	LatestLogByEventSigWithConfs(ctx context.Context, eventSig common.Hash, address common.Address, confs evmtypes.Confirmations) (*logpoller.Log, error)
	Replay(ctx context.Context, fromBlock int64) error
}

type lloConfigPoller struct {
	services.Service
	eng *services.Engine

	// runReplay       bool
	// fromBlock       int64
	// logPollInterval time.Duration

	// eventName string
	// eventSig  common.Hash
	// abi       *abi.ABI

	// lp       LLOLogPoller
	addr     common.Address
	contract *capabilities_registry.CapabilitiesRegistry

	donID        uint32
	capabilityID common.Hash
}

type CPConfig struct {
	Logger logger.Logger
	Client client.Client
	// LogPoller logpoller.LogPoller

	DonID                       uint32
	CapabilitiesRegistryAddress common.Address
	CapabilityName              string
	CapabilityVersion           string
}

type ConfigPoller interface {
	ocrtypes.ContractConfigTracker
	services.Service
}

func NewConfigPoller(cfg CPConfig) (ConfigPoller, error) {
	return newConfigPoller(cfg)
}

func newConfigPoller(cfg CPConfig) (*lloConfigPoller, error) {
	capabilitiesRegistryContract, err := capabilities_registry.NewCapabilitiesRegistry(cfg.CapabilitiesRegistryAddress, cfg.Client)
	if err != nil {
		return nil, err
	}

	cp := &lloConfigPoller{
		addr:         cfg.CapabilitiesRegistryAddress,
		contract:     capabilitiesRegistryContract,
		donID:        cfg.DonID,
		capabilityID: getHashedCapabilityId(cfg.CapabilityName, cfg.CapabilityVersion),
	}

	cp.Service, cp.eng = services.Config{
		Name:           fmt.Sprintf("LLOConfigPoller.%s", cfg.CapabilitiesRegistryAddress),
		NewSubServices: nil,
		Start:          cp.start,
		Close:          cp.close,
	}.NewServiceEngine(cfg.Logger)

	return cp, nil
}

func (c *lloConfigPoller) start(ctx context.Context) error {
	return nil
}

func (c *lloConfigPoller) close() error {
	return nil
}

func (c *lloConfigPoller) callGetDON(ctx context.Context) (donInfo capabilities_registry.CapabilitiesRegistryDONInfo, err error) {
	donInfo, err = c.contract.GetDON(&bind.CallOpts{
		// Pending     bool            // Whether to operate on the pending state or the last known one
		// From        common.Address  // Optional the sender address, otherwise the first account is used
		// BlockNumber *big.Int        // Optional the block number on which the call should be performed
		// BlockHash   common.Hash     // Optional the block hash on which the call should be performed
		// Context     context.Context // Network context to support cancellation and timeouts (nil = no timeout)
		Context: ctx,
	}, c.donID)
	if err != nil {
		failedRPCContractCalls.WithLabelValues(cp.client.ConfiguredChainID().String(), cp.aggregatorContractAddr.Hex()).Inc()
	}
	return
}

func (c *lloConfigPoller) callGetCapabilityConfigs(ctx context.Context) (donCapabilityConfig, globalCapabilityConfig []byte, err error) {
	donCapabilityConfig, globalCapabilityConfig, err = c.contract.GetCapabilityConfigs(&bind.CallOpts{
		// Pending     bool            // Whether to operate on the pending state or the last known one
		// From        common.Address  // Optional the sender address, otherwise the first account is used
		// BlockNumber *big.Int        // Optional the block number on which the call should be performed
		// BlockHash   common.Hash     // Optional the block hash on which the call should be performed
		// Context     context.Context // Network context to support cancellation and timeouts (nil = no timeout)
		Context: ctx,
	})
	if err != nil {
		failedRPCContractCalls.WithLabelValues(cp.client.ConfiguredChainID().String(), cp.aggregatorContractAddr.Hex()).Inc()
	}
	return
}

// getHashedCapabilityId copies the abi encoding of the capabilityId function
// in the contract, but for off-chain use
//
// See: https://github.com/smartcontractkit/chainlink/blob/e3a8ffead603d8cc87d97310ee284c80d825c7ab/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L757C1-L764C4
func getHashedCapabilityId(labelledName, version string) common.Hash {
	arguments := abi.Arguments{
		{
			Type: abi.StringTy,
		},
		{
			Type: abi.StringTy,
		},
	}

	// Encode the arguments using the Pack method
	data, err := arguments.Pack(labelledName, version)
	if err != nil {
		return nil, err
	}
	return utils.Keccak256Fixed(data)
}

// LatestBlockHeight implements types.ContractConfigTracker.
func (c *lloConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}

// LatestConfig implements types.ContractConfigTracker.
func (c *lloConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (contractConfig ocr2types.ContractConfig, err error) {
	return
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (c *lloConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocr2types.ConfigDigest, err error) {
	donInfo, err := c.callGetDON(ctx)
	if err != nil {
		return 0, configDigest, err
	}
	donCapabilityConfig, globalCapabilityConfig, err := c.callGetCapabilityConfigs(ctx)
	if err != nil {
		return 0, configDigest, err
	}

	// donInfo
	// Id                       uint32
	// ConfigCount              uint32
	// F                        uint8
	// IsPublic                 bool
	// AcceptsWorkflows         bool
	// NodeP2PIds               [][32]byte
	// CapabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration

	return
}

// Notify implements types.ContractConfigTracker.
func (c *lloConfigPoller) Notify() <-chan struct{} {
	return nil
}

// func (c *lloConfigPoller) Replay(ctx context.Context, fromBlock int64) error {
//     return c.lp.Replay(ctx, fromBlock)
// }

// func (c *lloConfigPoller) Ready() error {
//     return nil
// }

// func (c *lloConfigPoller) HealthReport() map[string]error {
//     return make(map[string]error)
// }

// func (c *lloConfigPoller) Name() string {
//     return c.Name()
// }

// func (c *lloConfigPoller) contractConfig() evmtypes.ContractConfig {
//     return evmtypes.ContractConfig{
// ConfigDigest:          c.cfg.ConfigDigest,
// ConfigCount:           c.cfg.ConfigCount,
// Signers:               toOnchainPublicKeys(c.cfg.Config.Signers),
// Transmitters:          toOCRAccounts(c.cfg.Config.Transmitters),
// F:                     c.cfg.Config.F,
// OnchainConfig:         []byte{},
// OffchainConfigVersion: c.cfg.Config.OffchainConfigVersion,
// OffchainConfig:        c.cfg.Config.OffchainConfig,
// }
// }

// PublicConfig returns the OCR configuration as a PublicConfig so that we can
// access ReportingPluginConfig and other fields prior to launching the plugins.
// func (c *lloConfigPoller) PublicConfig() (ocr3confighelper.PublicConfig, error) {
//     return ocr3confighelper.PublicConfigFromContractConfig(false, c.contractConfig())
// }

// func toOnchainPublicKeys(signers [][]byte) []types.OnchainPublicKey {
//     keys := make([]types.OnchainPublicKey, len(signers))
//     for i, signer := range signers {
//         keys[i] = types.OnchainPublicKey(signer)
//     }
//     return keys
// }

// func toOCRAccounts(transmitters [][]byte) []types.Account {
//     accounts := make([]types.Account, len(transmitters))
//     for i, transmitter := range transmitters {
//         // TODO: string-encode the transmitter appropriately to the dest chain family.
//         accounts[i] = types.Account(gethcommon.BytesToAddress(transmitter).Hex())
//     }
//     return accounts
// }

// var _ types.ContractConfigTracker = (*configTracker)(nil)

// var transmitAccounts []ocrtypes.Account
// for _, addr := range unpacked.OffchainTransmitters {
//     transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr)))
// }
// var signers []ocrtypes.OnchainPublicKey
// for _, addr := range unpacked.Signers {
//     addr := addr
//     signers = append(signers, addr[:])
// }

// return ocrtypes.ContractConfig{
//     ConfigDigest:          unpacked.ConfigDigest,
//     ConfigCount:           unpacked.ConfigCount,
//     Signers:               signers,
//     Transmitters:          transmitAccounts,
//     F:                     unpacked.F,
//     OnchainConfig:         unpacked.OnchainConfig,
//     OffchainConfigVersion: unpacked.OffchainConfigVersion,
//     OffchainConfig:        unpacked.OffchainConfig,
// }, nil
