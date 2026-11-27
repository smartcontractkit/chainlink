package aptos

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	pkgerrors "github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func donOraclePublicKeys(ctx context.Context, don *cre.Don) ([][]byte, error) {
	workers, err := don.Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to list worker nodes for DON %q: %w", don.Name, err)
	}

	oracles := make([][]byte, 0, len(workers))
	for _, worker := range workers {
		ocr2ID := ""
		if worker.Keys != nil && worker.Keys.OCR2BundleIDs != nil {
			ocr2ID = worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos]
		}
		if ocr2ID == "" {
			fetchedID, err := worker.Clients.GQLClient.FetchOCR2KeyBundleID(ctx, strings.ToUpper(chainselectors.FamilyAptos))
			if err != nil {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q and fallback fetch failed: %w", worker.Name, don.Name, err)
			}
			if fetchedID == "" {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q", worker.Name, don.Name)
			}
			ocr2ID = fetchedID
			if worker.Keys != nil {
				if worker.Keys.OCR2BundleIDs == nil {
					worker.Keys.OCR2BundleIDs = make(map[string]string)
				}
				worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos] = ocr2ID
			}
		}

		exported, err := worker.ExportOCR2Keys(ocr2ID)
		if err != nil {
			return nil, fmt.Errorf("failed to export Aptos OCR2 key for worker %q (bundle %s): %w", worker.Name, ocr2ID, err)
		}
		pubkey, err := parseOCR2OnchainPublicKey(exported.OnchainPublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid Aptos OCR2 onchain public key for worker %q: %w", worker.Name, err)
		}
		oracles = append(oracles, pubkey)
	}

	return oracles, nil
}

func p2pToTransmitterMapForWorkers(workers []*cre.NodeMetadata) (map[string]string, error) {
	if len(workers) == 0 {
		return nil, pkgerrors.New("no DON worker nodes provided")
	}

	p2pToTransmitterMap := make(map[string]string)
	for _, worker := range workers {
		if worker.Keys == nil || worker.Keys.P2PKey == nil {
			return nil, fmt.Errorf("missing P2P key for worker index %d", worker.Index)
		}

		account := ""
		if worker.Keys.Aptos != nil {
			account = worker.Keys.Aptos.Account
		}
		if account == "" {
			return nil, fmt.Errorf("missing Aptos account for worker index %d", worker.Index)
		}

		transmitter, err := normalizeTransmitter(account)
		if err != nil {
			return nil, fmt.Errorf("invalid Aptos transmitter for worker index %d: %w", worker.Index, err)
		}

		peerKey := hex.EncodeToString(worker.Keys.P2PKey.PeerID[:])
		p2pToTransmitterMap[peerKey] = transmitter
	}

	if len(p2pToTransmitterMap) == 0 {
		return nil, pkgerrors.New("no Aptos transmitters found for DON workers")
	}

	return p2pToTransmitterMap, nil
}

func aptosDonIDUint32(donID uint64) (uint32, error) {
	if donID > uint64(^uint32(0)) {
		return 0, fmt.Errorf("don id %d exceeds u32", donID)
	}
	return uint32(donID), nil
}

func parseOCR2OnchainPublicKey(hexValue string) ([]byte, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(hexValue), "ocr2on_aptos_")
	decoded, err := hex.DecodeString(trimmed)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
