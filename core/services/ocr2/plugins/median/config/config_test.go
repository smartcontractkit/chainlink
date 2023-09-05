package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePluginConfig(t *testing.T) {
	for _, s := range []struct {
		name     string
		pipeline string
	}{
		{"empty", ""},
		{"blank", " "},
		{"foo", "foo"},
	} {
		t.Run(s.name, func(t *testing.T) {
			config := PluginConfig{JuelsPerFeeCoinPipeline: s.pipeline}
			assert.Error(t, config.ValidatePluginConfig())
		})
	}
}
