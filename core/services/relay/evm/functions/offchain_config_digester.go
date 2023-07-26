package functions

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	_                     types.OffchainConfigDigester = FunctionsOffchainConfigDigester{}
	FunctionsDigestPrefix                              = types.ConfigDigestPrefixEVM
	// In order to support multiple OCR plugins with a single jobspec & OCR2Base contract, each plugin must have a unique config digest.
	// This is accomplished by overriding the single config digest from the contract with a unique prefix for each plugin via this custom offchain digester & config poller.
	ThresholdDigestPrefix = types.ConfigDigestPrefix(7)
	S4DigestPrefix        = types.ConfigDigestPrefix(8)
)

type FunctionsOffchainConfigDigester struct {
	PluginType   FunctionsPluginType
	BaseDigester evmutil.EVMOffchainConfigDigester
}

func (d FunctionsOffchainConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	configDigest, err := d.BaseDigester.ConfigDigest(cc)
	if err != nil {
		return types.ConfigDigest{}, err
	}

	var prefix types.ConfigDigestPrefix
	switch d.PluginType {
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

func (d FunctionsOffchainConfigDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	switch d.PluginType {
	case FunctionsPlugin:
		return FunctionsDigestPrefix, nil
	case ThresholdPlugin:
		return ThresholdDigestPrefix, nil
	case S4Plugin:
		return S4DigestPrefix, nil
	default:
		return 0, fmt.Errorf("unknown plugin type: %v", d.PluginType)
	}
}
