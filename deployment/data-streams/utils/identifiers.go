package utils

import (
	"fmt"
	"regexp"
)

const (
	ProductLabel = "data-streams"
)

// DonIdentifier generates a unique identifier for a DON based on its ID and name.
// All non-alphanumeric characters are replaced with underscores due to the limiting requirements of
// Job Distributor label keys.
func DonIdentifier(donID uint64, donName string) string {
	cleanDONName := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(donName, "_")
	return fmt.Sprintf("don-%d-%s", donID, cleanDONName)
}
