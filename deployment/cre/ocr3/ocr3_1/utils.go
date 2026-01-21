package ocr3_1

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func HexStringTo32ByteArray(s string) (*[32]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}
	if len(b) != 32 {
		return nil, fmt.Errorf("hex string must be 32 bytes, got %d", len(b))
	}
	b32 := [32]byte(b)
	return &b32, nil
}

func VerifyAndExtractOCR3_1Fields(prevConfigDigestStr string, prevSeqNr uint64, prevHistoryDigestStr string) (*types.ConfigDigest, *types.HistoryDigest, error) {
	arePreviousFieldsSet := prevConfigDigestStr != ""
	if arePreviousFieldsSet && (prevSeqNr == uint64(0) || prevHistoryDigestStr == "") {
		return nil, nil, errors.New("PrevConfigDigest, PrevSeqNr, and PrevHistoryDigest must all be set or all be nil")
	}
	if !arePreviousFieldsSet && (prevSeqNr != uint64(0) || prevHistoryDigestStr != "") {
		return nil, nil, errors.New("PrevConfigDigest, PrevSeqNr, and PrevHistoryDigest must all be set or all be nil")
	}
	var prevConfigDigest *types.ConfigDigest
	if prevConfigDigestStr != "" {
		bytes32, err := HexStringTo32ByteArray(prevConfigDigestStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid OracleConfig.PrevConfigDigest, should be hex encoded 32 bytes: %w", err)
		}
		prevConfigDigest = (*types.ConfigDigest)(bytes32)
	}
	var prevHistoryDigest *types.HistoryDigest
	if prevHistoryDigestStr != "" {
		bytes32, err := HexStringTo32ByteArray(prevHistoryDigestStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid OracleConfig.PrevHistoryDigest, should be hex encoded 32 bytes: %w", err)
		}
		prevHistoryDigest = (*types.HistoryDigest)(bytes32)
	}
	return prevConfigDigest, prevHistoryDigest, nil
}
