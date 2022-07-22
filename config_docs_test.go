package chainlink

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/internal/config"
)

//go:embed docs/CONFIG.md
var markdown string

func TestConfigDocs(t *testing.T) {
	got, err := config.GenerateDocs()
	assert.NoError(t, err, "invalid config docs")
	assert.Equal(t, markdown, got, "docs/CONFIG.md is out of date. Run 'make config-docs' to regenerate.")
    

}
