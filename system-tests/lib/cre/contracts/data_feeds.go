package contracts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// GetDecimalsFromFeedID extracts the number of decimals from a feed ID
// Feed ID format is assumed to be a hex string (with or without 0x prefix)
// The 7th byte (index 6 when zero-based) determines the decimal count
func GetDecimalsFromFeedID(feedID string) (int, error) {
	// Remove 0x prefix if present
	cleanFeedID := strings.TrimPrefix(feedID, "0x")

	// Ensure the feed ID is long enough
	if len(cleanFeedID) < 14 { // Need at least 7 bytes (14 hex chars)
		return 0, fmt.Errorf("feed ID too short: %s", feedID)
	}

	// Extract the 7th byte (bytes are represented by 2 hex characters)
	// Position 12-13 represents the 7th byte (0-indexed)
	dataTypeByte := cleanFeedID[12:14]

	// Convert the hex string to a byte
	dataTypeInt, err := strconv.ParseUint(dataTypeByte, 16, 8)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse data type byte")
	}

	// Map the data type byte to decimals
	switch {
	case dataTypeInt == 0x20:
		return 0, nil // Integer (Decimal0)
	case dataTypeInt >= 0x21 && dataTypeInt <= 0x60:
		return int(dataTypeInt - 0x20), nil // DecimalN where N is (dataTypeInt - 0x20)
	default:
		return 0, fmt.Errorf("unknown data type: 0x%s", dataTypeByte)
	}
}
