package cre

// NewUncheckedDonsMetadata builds DonsMetadata for unit tests without NewDonsMetadata validation.
//
// Pairing and deploy resolver tests need minimal DON lists (names, flags, don_family) that would
// fail or be awkward to construct through the full nodeset → DonMetadata constructor path.
func NewUncheckedDonsMetadata(dons []*DonMetadata) *DonsMetadata {
	return &DonsMetadata{dons: dons}
}
