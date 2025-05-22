package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	ProductLabel = "data-streams"
)

// DonIDLabel generates a unique identifier for a DON based on its ID and name.
// All non-alphanumeric characters are replaced with underscores due to the limiting requirements of
// Job Distributor label keys.
func DonIDLabel(donID uint64, donName string) string {
	cleanDONName := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(donName, "_")
	return fmt.Sprintf("don-%d-%s", donID, cleanDONName)
}

func StreamIDLabel(streamID uint32) string {
	return fmt.Sprintf("stream-id-%d", streamID)
}

func StreamIDFromLabel(streamIDLabel string) (uint32, error) {
	matches := regexp.MustCompile(`stream-id-([0-9]+)`).FindStringSubmatch(streamIDLabel)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid stream ID label: %s", streamIDLabel)
	}
	streamID, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse stream ID: %w", err)
	}
	return uint32(streamID), nil
}
