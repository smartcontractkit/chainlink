package engine

// Product bundles everything product-specific the changelog engine needs:
// which repositories to track, how to label reports, and how to normalize
// product-specific ref shapes (e.g. image tags) into git refs.
//
// The engine itself is product-agnostic: a Product is passed in by the
// caller. Product definitions live in internal/products/<name>/ — see
// internal/products/ccip for the reference implementation.
//
// To add support for another product (e.g. Core releases):
//  1. Create internal/products/<name>/ with a Product value (copy ccip).
//  2. Wire it in cmd/release-changelog/main.go — or, once a second
//     product actually exists, add a --product CLI flag / workflow input
//     that selects among them.
type Product struct {
	// DisplayName labels report headers and Slack messages (e.g. "CCIP").
	DisplayName string
	// Repos lists the tracked repositories; see RepoConfig docs.
	Repos []RepoConfig
	// NormalizeRef maps product-specific ref shapes (image tags, image URIs)
	// to git refs before resolution. Nil means "treat all inputs as git
	// refs already".
	NormalizeRef func(string) string
}

// normalize maps ref through the product's normalizer, if any.
func (p Product) normalize(ref string) string {
	if p.NormalizeRef == nil {
		return ref
	}
	return p.NormalizeRef(ref)
}
