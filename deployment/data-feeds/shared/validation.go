package shared

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// spec- https://docs.google.com/document/d/13ciwTx8lSUfyz1IdETwpxlIVSn1lwYzGtzOBBTpl5Vg/edit?tab=t.0#heading=h.dxx2wwn1dmoz
func ValidateFeedID(feedID string) error {
	// Check for "0x" prefix and remove it
	feedID = strings.TrimPrefix(feedID, "0x")

	if len(feedID) != 32 {
		return fmt.Errorf("invalid feed ID length. Expected 32 characters, got %d", len(feedID))
	}

	bytes, err := hex.DecodeString(feedID)
	if err != nil {
		return fmt.Errorf("invalid feed ID format: %w", err)
	}

	// Validate format byte
	format := bytes[0]
	if format != 0x01 && format != 0x02 {
		return errors.New("invalid format byte")
	}

	// bytes (1-4) are random bytes, so no validation needed

	// Validate attribute bucket bytes (5-6)
	attributeBucket := [2]byte{bytes[5], bytes[6]}
	attributeBucketHex := hex.EncodeToString(attributeBucket[:])
	if attributeBucketHex != "0003" && attributeBucketHex != "0700" {
		return errors.New("invalid attribute bucket bytes")
	}

	// Validate data type byte (7)
	dataType := bytes[7]
	validDataTypes := map[byte]bool{
		0x00: true, 0x01: true, 0x02: true, 0x03: true, 0x04: true,
	}
	for i := 0x20; i <= 0x60; i++ {
		validDataTypes[byte(i)] = true
	}
	if !validDataTypes[dataType] {
		return errors.New("invalid data type byte")
	}

	// Validate reserved bytes (8-15) are zero
	for i := 8; i < 16; i++ {
		if bytes[i] != 0x00 {
			return errors.New("reserved bytes must be zero")
		}
	}

	return nil
}
