package config

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
)

func TestDoc(t *testing.T) {
	var c chainlink.Config
	d := toml.NewDecoder(strings.NewReader(docsTOML))
	// Note: using v1 of go-toml since v2 provides no feedback about which keys
	d.Strict(true) // Ensure no extra fields
	err := d.Decode(&c)
	if err != nil && strings.Contains(err.Error(), "undecoded keys: ") {
		t.Errorf("Docs contain extra fields: %v", err)
	} else {
		require.NoError(t, err)
	}

	cfgtest.AssertFieldsNotNil(t, c)

	//TODO validate defaults? ensure non-zero examples?
}
