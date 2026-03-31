package types

import (
	"errors"
	"fmt"
)

var (
	ErrGlobalWorkflowCountLimitReached   = errors.New("global workflow count limit reached")
	ErrPerOwnerWorkflowCountLimitReached = errors.New("per owner workflow count limit reached")
)

// ArtifactFetchError represents an internal failure to fetch a workflow artifact.
// It preserves full details for developer debugging while providing a deterministic
// customer-facing message suitable for aggregation across nodes.
type ArtifactFetchError struct {
	ArtifactType string // "binary" or "config"
	URL          string
	Err          error
}

func (e *ArtifactFetchError) Error() string {
	return fmt.Sprintf("failed to fetch %s from %s : %s", e.ArtifactType, e.URL, e.Err)
}

func (e *ArtifactFetchError) Unwrap() error {
	return e.Err
}

func (e *ArtifactFetchError) CustomerError() string {
	return fmt.Sprintf("Internal error: failed to fetch workflow %s from storage. Contact support if this persists.", e.ArtifactType)
}
