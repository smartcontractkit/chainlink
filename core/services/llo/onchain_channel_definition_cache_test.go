package llo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

func Test_ChannelDefinitionCache(t *testing.T) {
	t.Run("Definitions", func(t *testing.T) {
		// NOTE: this is covered more thoroughly in the integration tests
		dfns := llotypes.ChannelDefinitions(map[llotypes.ChannelID]llotypes.ChannelDefinition{
			1: {
				ReportFormat:  llotypes.ReportFormat(43),
				ChainSelector: 42,
				StreamIDs:     []llotypes.StreamID{1, 2, 3},
			},
		})

		cdc := &channelDefinitionCache{definitions: dfns}

		assert.Equal(t, dfns, cdc.Definitions())
	})
}
