package chainlink

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/config/v2/docs"
)

//go:embed docs/CONFIG.md
var markdown string

func TestConfigDocs(t *testing.T) {
	got, err := docs.GenerateDocs()
	assert.NoError(t, err, "invalid config docs")
	assert.Equal(t, markdown, got, "docs/CONFIG.md is out of date. Run 'make config-docs' to regenerate.")

}
