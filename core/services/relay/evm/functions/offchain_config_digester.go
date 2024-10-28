package functions

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	_                     types.OffchainConfigDigester        = &functionsOffchainConfigDigester{}
	_                     evmRelayTypes.RouteUpdateSubscriber = &functionsOffchainConfigDigester{}
	FunctionsDigestPrefix                                     = types.ConfigDigestPrefixEVMSimple
	// In order to support multiple OCR plugins with a single jobspec & OCR2Base contract, each plugin must have a unique config digest.
	// This is accomplished by overriding the single config digest from the contract with a unique prefix for each plugin via this custom offchain digester & config poller.
	ThresholdDigestPrefix = types.ConfigDigestPrefix(7)
	S4DigestPrefix        = types.ConfigDigestPrefix(8)
)

type functionsOffchainConfigDigester struct {
	pluginType      FunctionsPluginType
	chainID         uint64
	contractAddress atomic.Pointer[common.Address]
}

func NewFunctionsOffchainConfigDigester(pluginType FunctionsPluginType, chainID uint64) *functionsOffchainConfigDigester {
	return &functionsOffchainConfigDigester{
		pluginType: pluginType,
		chainID:    chainID,
	}
}

func (d *functionsOffchainConfigDigester) ConfigDigest(ctx context.Context, cc types.ContractConfig) (types.ConfigDigest, error) {
	contractAddress := d.contractAddress.Load()
	if contractAddress == nil {
		return types.ConfigDigest{}, errors.New("contract address not set")
	}
	baseDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         d.chainID,
		ContractAddress: *contractAddress,
	}

	configDigest, err := baseDigester.ConfigDigest(ctx, cc)
	if err != nil {
		return types.ConfigDigest{}, err
	}

	var prefix types.ConfigDigestPrefix
	switch d.pluginType {
	case FunctionsPlugin:
		prefix = FunctionsDigestPrefix
	case ThresholdPlugin:
		prefix = ThresholdDigestPrefix
	case S4Plugin:
		prefix = S4DigestPrefix
	default:
		return types.ConfigDigest{}, errors.New("unknown plugin type")
	}

	binary.BigEndian.PutUint16(configDigest[:2], uint16(prefix))

	return configDigest, nil
}

func (d *functionsOffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (types.ConfigDigestPrefix, error) {
	switch d.pluginType {
	case FunctionsPlugin:
		return FunctionsDigestPrefix, nil
	case ThresholdPlugin:
		return ThresholdDigestPrefix, nil
	case S4Plugin:
		return S4DigestPrefix, nil
	default:
		return 0, fmt.Errorf("unknown plugin type: %v", d.pluginType)
	}
}

// called from LogPollerWrapper in a separate goroutine
func (d *functionsOffchainConfigDigester) UpdateRoutes(ctx context.Context, activeCoordinator common.Address, proposedCoordinator common.Address) error {
	d.contractAddress.Store(&activeCoordinator)
	return nil
}
