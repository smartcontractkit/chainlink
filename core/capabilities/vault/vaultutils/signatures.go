package vaultutils

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func ValidateSignatures(resp *vaulttypes.SignedOCRResponse, allowedSigners []common.Address, minRequired int) error {
	if len(resp.Context) < 64 {
		return fmt.Errorf("context too short: expected min 64 bytes, got %d bytes", len(resp.Context))
	}

	if len(resp.Signatures) < minRequired {
		return fmt.Errorf("not enough signatures: expected min %d, got %d", minRequired, len(resp.Signatures))
	}

	// The context contains:
	// 0:32 -> config digest
	// 32:64 -> epoch + round, namely:
	//   - 0:27 -> zero padding
	//   - 27:31 -> sequence number (big endian uint32)
	//   - 31:32 -> zero round value
	// 64:96 -> extra hash (not used by the vault plugin)
	cd, epochRound := resp.Context[:32], resp.Context[32:64]
	configDigest, err := ocr2types.BytesToConfigDigest(cd)
	if err != nil {
		return fmt.Errorf("invalid config digest in signature: %w", err)
	}

	epoch := binary.BigEndian.Uint32(epochRound[27:31])
	round := epochRound[31]

	fullHash := ocr2key.ReportToSigData(ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}, []byte(resp.Payload))

	validSigners := map[common.Address]bool{}
	for _, s := range resp.Signatures {
		signerPubkey, err := crypto.SigToPub(fullHash, s)
		if err != nil {
			return fmt.Errorf("invalid signature: %w", err)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)

		for _, as := range allowedSigners {
			if as.Hex() == signerAddr.Hex() {
				validSigners[signerAddr] = true
				break
			}
		}

		if len(validSigners) >= minRequired {
			return nil
		}
	}

	return fmt.Errorf("only %d valid signatures, need at least %d", len(validSigners), minRequired)
}
