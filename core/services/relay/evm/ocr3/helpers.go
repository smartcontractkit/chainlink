package ocr3

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func configTrackerFilterName(id relay.ID, addr common.Address) string {
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
func TransmitterCombiner(masterConfig ocrtypes.ContractConfig, followerConfigs []ocrtypes.ContractConfig) (ocrtypes.ContractConfig, error) {
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

	var combinedTransmitters []ocrtypes.Account
	for i, signer := range masterConfig.Signers {
		// the transmitter index is the same as the signer index for the same config object.
		// this is enforced in the standard OCR3Base.setOCR3Config method.
		transmitters := []string{string(masterConfig.Transmitters[i])}

		for _, followerConfig := range followerConfigs {
			// signer might be at a different index than master chain (but ideally shouldn't be)
			// so we can't just use i here.
			signerIdx := slices.IndexFunc(followerConfig.Signers, func(opk ocrtypes.OnchainPublicKey) bool {
				return hexutil.Encode(opk) == hexutil.Encode(signer)
			})
			if signerIdx == -1 {
				// signer not found, bad config
				return ocrtypes.ContractConfig{}, fmt.Errorf("unable to find signer %x (oracle index %d) in follower config %+v", signer, i, followerConfig)
			}
			// the transmitter index is the same as the signer index for the same config object.
			transmitters = append(transmitters, string(followerConfig.Transmitters[signerIdx]))
		}
		combinedTransmitter := joinTransmitters(transmitters)
		combinedTransmitters = append(combinedTransmitters, ocrtypes.Account(combinedTransmitter))
	}

	// sanity check
	if len(combinedTransmitters) != len(masterConfig.Signers) {
		return ocrtypes.ContractConfig{}, fmt.Errorf("unexpected length mismatch between combined transmitters (%d) and master config signers (%d)", len(combinedTransmitters), len(masterConfig.Signers))
	}

	toReturn.Transmitters = combinedTransmitters
	return toReturn, nil
}

// joinTransmitters is a helper that combines many transmitters into one
// Note that this is pulled out so that it can be used in the CombinerFn
// and the contract transmitter since the output of FromAccount() in the
// ContractTransmitter and the ContractConfig.Transmitters output for a
// particular signer must match in order for OCR3 to work.
func joinTransmitters(transmitters []string) string {
	// sort first to ensure deterministic ordering
	slices.Sort(transmitters)
	return strings.Join(transmitters, ",")
}
