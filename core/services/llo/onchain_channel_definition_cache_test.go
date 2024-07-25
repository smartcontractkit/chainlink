package llo

import (
	"testing"
)

func Test_ChannelDefinitionCache(t *testing.T) {
	t.Skip("waiting on https://github.com/smartcontractkit/chainlink/pull/13780")
	// t.Run("Definitions", func(t *testing.T) {
	//     // NOTE: this is covered more thoroughly in the integration tests
	//     dfns := llotypes.ChannelDefinitions(map[llotypes.ChannelID]llotypes.ChannelDefinition{
	//         1: {
	//             ReportFormat:  llotypes.ReportFormat(43),
	//             ChainSelector: 42,
	//             StreamIDs:     []llotypes.StreamID{1, 2, 3},
	//         },
	//     })

	//     cdc := &channelDefinitionCache{definitions: dfns}

	//     assert.Equal(t, dfns, cdc.Definitions())
	// })
}
