package metadata

// SimpleEnv represents the environment metadata for the exemplar domain.
type SimpleEnv struct {
	// DeployCounts is a map of chain selector to the number of contracts that have been deployed on that chain.
	DeployCounts map[uint64]int64 `json:"counts"`
}
