package vaultutils

import "strings"

// NormalizeOwner lowercases an Ethereum owner address for case-insensitive comparison.
// All comparison sites must use this function. When VaultOwnerAddressCanonicalizationEnabled
// is introduced, normalization at ingress will supersede comparison-site calls here.
func NormalizeOwner(owner string) string {
	return strings.ToLower(strings.TrimPrefix(owner, "0x"))
}
