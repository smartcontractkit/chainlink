package functions

import (
	"context"
	"database/sql"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ types.ConfigPoller = &configPoller{}

type FunctionsPluginType int

const (
	FunctionsPlugin FunctionsPluginType = iota
	ThresholdPlugin
	S4Plugin
)

type configPoller struct {
	lggr               logger.Logger
	filterName         string
	destChainLogPoller logpoller.LogPoller
	addr               common.Address
	pluginType         FunctionsPluginType
}

// ConfigSet Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/dev/OCR2Abstract.sol
var ConfigSet common.Hash

var defaultABI abi.ABI

const configSetEventName = "ConfigSet"

func init() {
	var err error
	abiPointer, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	defaultABI = *abiPointer
	ConfigSet = defaultABI.Events[configSetEventName].ID
}

func unpackLogData(d []byte) (*ocr2aggregator.OCR2AggregatorConfigSet, error) {
	unpacked := new(ocr2aggregator.OCR2AggregatorConfigSet)
	err := defaultABI.UnpackIntoInterface(unpacked, configSetEventName, d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack log data")
	}
	return unpacked, nil
}

func configFromLog(logData []byte, pluginType FunctionsPluginType) (ocrtypes.ContractConfig, error) {
	unpacked, err := unpackLogData(logData)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.String()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}

	// Replace the first two bytes of the config digest with the plugin type to avoid duplicate config digests between Functions plugins
	switch pluginType {
	case FunctionsPlugin:
		// FunctionsPluginType should already have the correct prefix, so this is a no-op
	case ThresholdPlugin:
		binary.BigEndian.PutUint16(unpacked.ConfigDigest[:2], uint16(ThresholdDigestPrefix))
	case S4Plugin:
		binary.BigEndian.PutUint16(unpacked.ConfigDigest[:2], uint16(S4DigestPrefix))
	default:
		return ocrtypes.ContractConfig{}, errors.New("unknown plugin type")
	}

	return ocrtypes.ContractConfig{
		ConfigDigest:          unpacked.ConfigDigest,
		ConfigCount:           unpacked.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     unpacked.F,
		OnchainConfig:         unpacked.OnchainConfig,
		OffchainConfigVersion: unpacked.OffchainConfigVersion,
		OffchainConfig:        unpacked.OffchainConfig,
	}, nil
}

func configPollerFilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2ConfigPoller", addr.String())
}

func NewFunctionsConfigPoller(pluginType FunctionsPluginType, destChainPoller logpoller.LogPoller, addr common.Address, lggr logger.Logger) (types.ConfigPoller, error) {
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configPollerFilterName(addr), EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	cp := &configPoller{
		lggr:               lggr,
		filterName:         configPollerFilterName(addr),
		destChainLogPoller: destChainPoller,
		addr:               addr,
		pluginType:         pluginType,
	}

	return cp, nil
}

func (cp *configPoller) Start() {}

func (cp *configPoller) Close() error {
	return nil
}

func (cp *configPoller) Notify() <-chan struct{} {
	return nil
}

func (cp *configPoller) Replay(ctx context.Context, fromBlock int64) error {
	return cp.destChainLogPoller.Replay(ctx, fromBlock)
}

func (cp *configPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latest, err := cp.destChainLogPoller.LatestLogByEventSigWithConfs(ConfigSet, cp.addr, 1, pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ocrtypes.ConfigDigest{}, nil
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	latestConfigSet, err := configFromLog(latest.Data, cp.pluginType)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

func (cp *configPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := cp.destChainLogPoller.Logs(int64(changedInBlock), int64(changedInBlock), ConfigSet, cp.addr, pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	latestConfigSet, err := configFromLog(lgs[len(lgs)-1].Data, cp.pluginType)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	cp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet, nil
}

func (cp *configPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := cp.destChainLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest), nil
}
