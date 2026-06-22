package cre

// NewUncheckedDonsMetadata constructs DonsMetadata for unit tests without running validation.
func NewUncheckedDonsMetadata(dons []*DonMetadata) *DonsMetadata {
	return &DonsMetadata{dons: dons}
}
