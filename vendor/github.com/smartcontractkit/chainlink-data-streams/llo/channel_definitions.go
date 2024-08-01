package llo

import (
	"fmt"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

func VerifyChannelDefinitions(channelDefs llotypes.ChannelDefinitions) error {
	if len(channelDefs) > MaxOutcomeChannelDefinitionsLength {
		return fmt.Errorf("too many channels, got: %d/%d", len(channelDefs), MaxOutcomeChannelDefinitionsLength)
	}
	uniqueStreamIDs := make(map[llotypes.StreamID]struct{}, len(channelDefs))
	for channelID, cd := range channelDefs {
		if len(cd.Streams) == 0 {
			return fmt.Errorf("ChannelDefinition with ID %d has no streams", channelID)
		}
		for _, strm := range cd.Streams {
			if strm.Aggregator == 0 {
				return fmt.Errorf("ChannelDefinition with ID %d has stream %d with zero aggregator (this may indicate an uninitialized struct)", channelID, strm.StreamID)
			}
			uniqueStreamIDs[strm.StreamID] = struct{}{}
		}
	}
	if len(uniqueStreamIDs) > MaxObservationStreamValuesLength {
		return fmt.Errorf("too many unique stream IDs, got: %d/%d", len(uniqueStreamIDs), MaxObservationStreamValuesLength)
	}
	return nil
}
