package utils

import (
	"fmt"
)

const (
	ProductLabel = "data-streams"
)

// DonIdentifier generates a unique identifier for a DON based on its ID and name.
func DonIdentifier(donID uint64, donName string) string {
	return fmt.Sprintf("don-%d-%s", donID, donName)
}

// DSLabel generates a label for the Data Streams product using the provided label.
func DSLabel(label string) string {
	return fmt.Sprintf("%s-%s", ProductLabel, label)
}
