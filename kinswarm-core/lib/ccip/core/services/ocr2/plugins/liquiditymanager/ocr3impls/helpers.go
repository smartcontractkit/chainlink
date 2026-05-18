package ocr3impls

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func configTrackerFilterName(id commontypes.RelayID, addr common.Address) string {
	return logpoller.FilterName("OCR3 MultichainConfigTracker", id.String(), addr.String())
}

func unpackLogData(d []byte) (*no_op_ocr3.NoOpOCR3ConfigSet, error) {
	unpacked := new(no_op_ocr3.NoOpOCR3ConfigSet)
	err := defaultABI.UnpackIntoInterface(unpacked, "ConfigSet", d)
	if err != nil {
		return nil, err
	}
	return unpacked, nil
}

func configFromLog(logData []byte) (ocrtypes.ContractConfig, error) {
	unpacked, err := unpackLogData(logData)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}

	var transmitAccounts []ocrtypes.Account
	for _, addr := range unpacked.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range unpacked.Signers {
		addr := addr
		signers = append(signers, addr[:])
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

// TransmitterCombiner is a CombinerFn that combines all transmitter addresses
// for the same signer on many different chains into a single string.
func TransmitterCombiner(masterChain commontypes.RelayID, contractConfigs map[commontypes.RelayID]ocrtypes.ContractConfig) (ocrtypes.ContractConfig, error) {
	masterConfig, ok := contractConfigs[masterChain]
	if !ok {
		return ocrtypes.ContractConfig{}, fmt.Errorf("unable to find master chain %s in contract configs", masterChain)
	}

	toReturn := ocrtypes.ContractConfig{
		ConfigDigest: masterConfig.ConfigDigest,
		ConfigCount:  masterConfig.ConfigCount,
		Signers:      masterConfig.Signers,
		// Transmitters:          []ocrtypes.Account{}, // will be filled below
		F:                     masterConfig.F,
		OnchainConfig:         masterConfig.OnchainConfig,
		OffchainConfigVersion: masterConfig.OffchainConfigVersion,
		OffchainConfig:        masterConfig.OffchainConfig,
	}

	// group the transmitters for each signer into a single types.Account object
	var combinedTransmitters []ocrtypes.Account
	for i, signer := range masterConfig.Signers {
		// the transmitter index is the same as the signer index for the same config object.
		// this is enforced in the standard OCR3Base.setOCR3Config method.
		var transmitters []string

		for relayID, contractConfig := range contractConfigs {
			// signer might be at a different index than master chain (but ideally shouldn't be)
			// so we can't just use i here.
			signerIdx := slices.IndexFunc(contractConfig.Signers, func(opk ocrtypes.OnchainPublicKey) bool {
				return hexutil.Encode(opk) == hexutil.Encode(signer)
			})
			if signerIdx == -1 {
				// signer not found, bad config
				// this means that a signer on the main chain was not set on one of the follower chains
				return ocrtypes.ContractConfig{}, fmt.Errorf("unable to find signer %x (oracle index %d) in follower config %+v",
					signer, i, contractConfig)
			}
			if signerIdx >= len(contractConfig.Transmitters) {
				// should be impossible, since setOCR3Config enforces that the lengths
				// of the signers and transmitters be equal
				return ocrtypes.ContractConfig{}, fmt.Errorf("signer index %d out of bounds for follower config %+v on chain %s",
					signerIdx, contractConfig, relayID.String())
			}
			// the transmitter index is the same as the signer index for the same config object.
			transmitters = append(transmitters, EncodeTransmitter(relayID, contractConfig.Transmitters[signerIdx]))
		}
		combinedTransmitter := JoinTransmitters(transmitters)
		combinedTransmitters = append(combinedTransmitters, ocrtypes.Account(combinedTransmitter))
	}

	// sanity check
	if len(combinedTransmitters) != len(masterConfig.Signers) {
		return ocrtypes.ContractConfig{}, fmt.Errorf("unexpected length mismatch between combined transmitters (%d) and master config signers (%d)", len(combinedTransmitters), len(masterConfig.Signers))
	}

	toReturn.Transmitters = combinedTransmitters
	return toReturn, nil
}

// EncodeTransmitter encodes the provided relay ID and transmitter
// into a single string the following way:
// "<relayID.ChainID>:<transmitter>"
func EncodeTransmitter(relayID commontypes.RelayID, transmitter ocrtypes.Account) string {
	return fmt.Sprintf("%s:%s", relayID.ChainID, transmitter)
}

// JoinTransmitters is a helper that combines many transmitters into one
// Note that this is pulled out so that it can be used in the CombinerFn
// and the contract transmitter since the output of FromAccount() in the
// ContractTransmitter and the ContractConfig.Transmitters output for a
// particular signer must match in order for OCR3 to work.
func JoinTransmitters(transmitters []string) string {
	// sort first to ensure deterministic ordering
	slices.Sort(transmitters)
	return strings.Join(transmitters, ",")
}

// SplitMultiTransmitter splits a multi-transmitter string
// into a map of chainID -> signerIdx -> transmitter
// This is so that we can verify the config digest offchain in the
// offchain config digester.
func SplitMultiTransmitter(multiTransmitter ocrtypes.Account) (map[commontypes.RelayID]ocrtypes.Account, error) {
	toReturn := map[commontypes.RelayID]ocrtypes.Account{}
	for _, chainAndTransmitter := range strings.Split(string(multiTransmitter), ",") {
		parts := strings.Split(chainAndTransmitter, ":")
		if len(parts) != 2 {
			// assertion
			return nil, fmt.Errorf("split on ':' must contain exactly 2 parts, got: %d (%+v)",
				len(parts), parts)
		}
		chainID := commontypes.NewRelayID(relay.NetworkEVM, parts[0])
		transmitter := ocrtypes.Account(parts[1])

		if _, ok := toReturn[chainID]; ok {
			return nil, fmt.Errorf("same chain id appearing multiple times in parts %+v (multi-transmitter %s): %s",
				parts, multiTransmitter, chainID.String())
		}
		toReturn[chainID] = transmitter
	}
	return toReturn, nil
}
