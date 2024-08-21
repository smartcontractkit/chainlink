package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ ocr2types.ContractConfigTracker = (*lloConfigPoller)(nil)

type LLOLogPoller interface {
	RegisterFilter(ctx context.Context, filter logpoller.Filter) error
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
	LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error)
	Replay(ctx context.Context, fromBlock int64) error
}

type lloConfigPoller struct {
	services.Service
	eng *services.Engine

	eventName string
	eventSig  common.Hash
	abi       *abi.ABI

	lp       LLOLogPoller
	addr     common.Address
	contract *capabilities_registry.CapabilitiesRegistry
}

func NewLLOConfigPoller(ctx context.Context, lggr logger.Logger, client client.Client, lp LLOLogPoller, capabilitiesRegistryAddr common.Address) (evmRelayTypes.ConfigPoller, error) {
	return newLLOConfigPoller(ctx, lggr, client, lp, capabilitiesRegistryAddr)
}

func newLLOConfigPoller(ctx context.Context, lggr logger.Logger, client client.Client, lp LLOLogPoller, capabilitiesRegistryAddr common.Address) (*lloConfigPoller, error) {
	abi, err := capabilities_registry.CapabilitiesRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	capabilitiesRegistryContract, err := capabilities_registry.NewCapabilitiesRegistry(capabilitiesRegistryAddr, client)
	if err != nil {
		return nil, err
	}

	const eventName = "ConfigSet"
	cp := &lloConfigPoller{
		eventName: eventName,
		eventSig:  abi.Events[eventName].ID,
		abi:       abi,
		lp:        lp,
		addr:      capabilitiesRegistryAddr,
		contract:  capabilitiesRegistryContract,
	}

	cp.Service, cp.eng = services.Config{
		Name:           fmt.Sprintf("LLOConfigPoller.%s", capabilitiesRegistryAddr),
		NewSubServices: nil,
		Start:          cp.start,
		Close:          cp.close,
	}.NewServiceEngine(lggr)

	return cp, nil
}

func lloConfigPollerFilterName(addr common.Address) string {
	return logpoller.FilterName("LLOConfigPoller", addr.String())
}

func (c *lloConfigPoller) decode(rawLog []byte) (*capabilities_registry.CapabilitiesRegistryConfigSet, error) {
	unpacked := new(capabilities_registry.CapabilitiesRegistryConfigSet)
	err := c.abi.UnpackIntoInterface(unpacked, c.eventName, rawLog)
	return unpacked, err
}

func (c *lloConfigPoller) start(ctx context.Context) error {
	if c.runReplay && c.fromBlock != 0 {
		// Only replay if it's a brand new job.
		c.eng.Go(func(ctx context.Context) {
			c.eng.Infow("starting replay for config", "fromBlock", c.fromBlock)
			if err := c.lp.Replay(ctx, int64(c.fromBlock)); err != nil {
				c.eng.Errorw("error replaying for config", "err", err)
			} else {
				c.eng.Infow("completed replaying for config", "fromBlock", c.fromBlock)
			}
		})
	}
	err := c.lp.RegisterFilter(ctx, logpoller.Filter{Name: configPollerFilterName(c.addr), EventSigs: []common.Hash{c.eventSig}, Addresses: []common.Address{c.addr}})
	if err != nil {
		return err
	}
	return nil
}

func (c *lloConfigPoller) close() error {
	return nil
}

// LatestBlockHeight implements types.ContractConfigTracker.
func (c *lloConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}

// LatestConfig implements types.ContractConfigTracker.
func (c *lloConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (contractConfig types.ContractConfig, err error) {
	return
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (c *lloConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest types.ConfigDigest, err error) {
	return
}

// Notify implements types.ContractConfigTracker.
func (c *lloConfigPoller) Notify() <-chan struct{} {
	return nil
}

func (c *lloConfigPoller) contractConfig() types.ContractConfig {
	return types.ContractConfig{
		// ConfigDigest:          c.cfg.ConfigDigest,
		// ConfigCount:           c.cfg.ConfigCount,
		// Signers:               toOnchainPublicKeys(c.cfg.Config.Signers),
		// Transmitters:          toOCRAccounts(c.cfg.Config.Transmitters),
		// F:                     c.cfg.Config.F,
		// OnchainConfig:         []byte{},
		// OffchainConfigVersion: c.cfg.Config.OffchainConfigVersion,
		// OffchainConfig:        c.cfg.Config.OffchainConfig,
	}
}

// PublicConfig returns the OCR configuration as a PublicConfig so that we can
// access ReportingPluginConfig and other fields prior to launching the plugins.
func (c *lloConfigPoller) PublicConfig() (ocr3confighelper.PublicConfig, error) {
	return ocr3confighelper.PublicConfigFromContractConfig(false, c.contractConfig())
}

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
