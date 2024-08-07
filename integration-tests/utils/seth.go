package utils

import (
	"fmt"
	"testing"

	pkg_seth "github.com/smartcontractkit/seth"
)

// DynamicArtifactDirConfigFn returns a function that sets Seth's artifacts directory to a unique directory for the test
func DynamicArtifactDirConfigFn(t *testing.T) func(*pkg_seth.Config) error {
	return func(cfg *pkg_seth.Config) error {
		cfg.ArtifactsDir = fmt.Sprintf("seth_artifacts/%s", t.Name())
		return nil
	}
}
