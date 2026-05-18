package ocr3impls

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocr2plustypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

var _ types.OffchainConfigDigester = (*MultichainConfigDigester)(nil)

// MultichainConfigDigest is an offchain config digester implementation that is
// aware of the ContractConfig that is constructed from configs on multiple chains.
type MultichainConfigDigester struct {
	MasterChainDigester evmutil.EVMOffchainConfigDigester
}

// ConfigDigest calculates the config digest of the master chain from the provided
// combined contract config.
// The combined contract config will join transmitters from many chains into a since ocrtypes.Account
// object.
// Therefore we need a special implementation in order to extract the transmitters for the
// master chain only in order to calculate the correct master chain config digest.
func (d MultichainConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	signers := []common.Address{}
	for i, signer := range cc.Signers {
		if len(signer) != 20 {
			return types.ConfigDigest{}, fmt.Errorf("%v-th evm signer should be a 20 byte address, but got %x", i, signer)
		}
		a := common.BytesToAddress(signer)
		signers = append(signers, a)
	}

	if len(signers) != len(cc.Transmitters) {
		return types.ConfigDigest{}, fmt.Errorf("number of signers (%v) does not match number of transmitters (%v)", len(signers), len(cc.Transmitters))
	}

	// assemble the transmitters for the master chain, since thats
	// the chain whose config digest we want to calculate
	var masterChainTransmitters []types.Account
	for i, multiTransmitter := range cc.Transmitters {
		split, err := SplitMultiTransmitter(multiTransmitter)
		if err != nil {
			return types.ConfigDigest{}, fmt.Errorf("unable to split multi-transmitter %s: %w", multiTransmitter, err)
		}
		masterRelayID := commontypes.NewRelayID(relay.NetworkEVM, strconv.FormatUint(d.MasterChainDigester.ChainID, 10))
		if _, ok := split[masterRelayID]; !ok {
			return types.ConfigDigest{}, fmt.Errorf("multi-transmitter %s does not contain a transmitter for master chain %s", multiTransmitter, masterRelayID)
		}
		transmitter := split[masterRelayID]
		if !strings.HasPrefix(string(transmitter), "0x") || len(transmitter) != 42 || !common.IsHexAddress(string(transmitter)) {
			return types.ConfigDigest{}, fmt.Errorf("%v-th evm transmitter should be a 42 character Ethereum address string, but got '%v'", i, transmitter)
		}
		masterChainTransmitters = append(masterChainTransmitters, transmitter)
	}

	if len(masterChainTransmitters) != len(signers) {
		return types.ConfigDigest{}, fmt.Errorf("number of signers (%v) does not match number of transmitters (%v)", len(signers), len(masterChainTransmitters))
	}

	return d.MasterChainDigester.ConfigDigest(types.ContractConfig{
		ConfigDigest:          cc.ConfigDigest,
		ConfigCount:           cc.ConfigCount,
		Signers:               cc.Signers,
		Transmitters:          masterChainTransmitters,
		F:                     cc.F,
		OnchainConfig:         cc.OnchainConfig,
		OffchainConfigVersion: cc.OffchainConfigVersion,
		OffchainConfig:        cc.OffchainConfig,
	})
}

func (d MultichainConfigDigester) ConfigDigestPrefix() (types.ConfigDigestPrefix, error) {
	return ocr2plustypes.ConfigDigestPrefixEVMSimple, nil
}
